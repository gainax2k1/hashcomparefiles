package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/gainax2k1/hashcomparefiles/internal/hashfile"
	"github.com/gainax2k1/hashcomparefiles/internal/logging"
	walkdir "github.com/gainax2k1/hashcomparefiles/internal/walkdir"
)

type Config struct {
	FilePath      string
	TrashPath     string
	TrashInfoPath string
	LogPath       string
	RemoveFlag    bool
	Minflag       int64
	Maxflag       int64
}

var maxFileSize int64 = math.MaxInt64

func main() {
	// Define flags and parse
	removeFlag := flag.Bool("remove", false, "Selectively choose which duplicates to trash or delete if desired")
	logFlag := flag.String("log", "slog.log", "Log filename, or 'default' for current directory log.log")
	minFlag := flag.Int64("min", 1, "Minimum filesize to include (in bytes")
	maxFlag := flag.Int64("max", maxFileSize, "Maximum filesize to include (in bytes)")
	verboseFlag := flag.Bool("v", false, "Output complete duplicate list to screen upon completion")

	flag.Parse()

	// Identify all paths to process (pipe or args)
	var targets []string

	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		// Data is being piped in
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			if path := strings.TrimSpace(scanner.Text()); path != "" {
				targets = append(targets, path)
			}
		}
	} else {
		// Use command line arguments if no pipe
		targets = flag.Args()
	}

	// Validate targets
	if len(targets) == 0 {
		fmt.Fprintf(os.Stderr, "Usage: %s [flags] <path>\n", os.Args[0])
		os.Exit(1)
	}

	// Create config struct with parsed values
	trashPath, trashInfoPath, err := configTrash()
	if err != nil {
		log.Fatalf("Error configuring trash: %v", err)

	}

	config := Config{
		TrashPath:     trashPath,
		TrashInfoPath: trashInfoPath,
		LogPath:       *logFlag,
		RemoveFlag:    *removeFlag,
		Minflag:       *minFlag,
		Maxflag:       *maxFlag,
	}
	if config.RemoveFlag {
		//Force verbose when doing removal to display submaps of duplicates
		*verboseFlag = true
	}

	// All output will be done through the logger, writing to file and/or screen based on config
	logger, err := logging.NewLogger(config.LogPath, *verboseFlag)
	if err != nil {
		log.Fatalf("Error creating logger: %v", err)
	}
	defer logger.Close() // ensure even if panic, logger successfully closes

	err = process(targets, config, logger)
	if err != nil {
		logger.Error("Error processing: %v", err)
	}
}

func process(targets []string, config Config, logger *logging.Logger) error {
	// 1. Make map (key=filesize, value=[]filepaths)
	// --- ignore all symlinks, zero size files
	// 2. For each key, if len(value) > 1, then run smaller hash on each file, make map of (key=hash, value=[]filepaths)
	// --- first pass, save on hashing on large files unless neccessary
	// 3. For each key, if len(value) > 1, then run full hash on each file, make map of (key=hash, value=[]filepaths)
	// --- 2nd pass, run fullhash on remaining files

	var spinnerCounter = 0
	// Get the number of available CPU cores and set max number of threads to that for optimal performance

	// FIRST PASS:
	fileSizeMap := make(map[int64][]string)
	totalCount := 0

	for _, path := range targets {
		// Map files by filesize
		dirMap, count, err := walkdir.WalkGetFileSizes(path, logger)
		if err != nil {
			logger.Error("Skipping %s due to error: %v", path, err)
			continue // Keep going with other targets!
		}

		totalCount += count

		for size, files := range dirMap {
			if size > config.Minflag && size < config.Maxflag {
				fileSizeMap[size] = append(fileSizeMap[size], files...)
			}
		}
	}

	logger.Info("pass_complete", "pass", 1, "count", totalCount)
	fmt.Printf("Filecount after pass (1/3): %d\n", totalCount)

	// SECOND PASS:
	totalCount = 0 //reset

	sem := make(chan struct{}, 32) // limit go routines to 128 for hashing, since hashing is CPU intensive, this should optimize performance without overwhelming the system. Adjust as needed based on testing and system capabilities.

	for filesize, files := range fileSizeMap {

		if len(files) == 1 { // only one file with this size, so unique
			continue // skip this file
		}
		// multiple files with this size, so we need to compare them
		for _, file := range files {
			wg.Add(1)
			sem <- struct{}{} // block if 128 goroutines are already running
			go func(f string, fs int64) {
				defer func() { <-sem }()
				processPartialFile(f, fs, logger)
			}(file, filesize)

			spinnerCounter++
			if spinnerCounter%100 == 0 {
				// Pulse/Spinner every 100 files to save CPU
				// \r clears the line, then we print the spinner and count
				fmt.Fprintf(os.Stderr, "\r %s Files processed: %d\r", getSpinner(spinnerCounter/100), spinnerCounter)

			}
			totalCount++
		}
	}
	wg.Wait() // Wait for all goroutines to finish before moving to the next step

	logger.Info("pass_complete", "pass", 2, "count", totalCount)
	fmt.Printf("Filecount after pass (2/3): %d\n", totalCount)
	spinnerCounter = 0

	fileSizeMap = nil // free memory from first pass, since we now have the partial hashes in memory, we can release the file size map

	// THIRD PASS:
	//finalDuplicates := make(map[string][]walkdir.FileInfo)
	totalCount = 0 //reset
	for smallHash, files := range firstPassHashes {
		if len(files) == 1 { // only one file with this size
			continue // skip this file, hash to be unique
		}
		for _, file := range files {
			if file.FileSize <= hashfile.PARTIALBYTELENGTH {
				mu.Lock()
				finalDuplicates[smallHash] = append(finalDuplicates[smallHash], file)
				mu.Unlock()
				totalCount++
				continue // use first hash, since file was already *fully* hashed
			}

			wg.Add(1)
			sem <- struct{}{} // block if 128 goroutines are already running
			go func(f string, fs int64) {
				defer func() { <-sem }()
				processFullFile(f, fs, logger)
			}(file.FilePath, file.FileSize)

			spinnerCounter++
			if spinnerCounter%100 == 0 {
				// Pulse/Spinner every 100 files to save CPU
				// \r clears the line, then we print the spinner and count
				fmt.Fprintf(os.Stderr, "\r %s Files processed: %d\r", getSpinner(spinnerCounter/100), spinnerCounter)
			}
			totalCount++
		}
	}
	wg.Wait()

	logger.Info("pass_complete", "pass", 3, "count", totalCount)
	fmt.Printf("Filecount after pass (3/3): %d\n", totalCount)

	firstPassHashes = nil // free memory from second pass, since we now have the full hashes in memory, we can release the partial hash map

	//shrink map
	finalMap, totalCount := filterDuplicates(finalDuplicates)

	logger.Info("pass_complete", "pass", 4, "count", totalCount)
	fmt.Printf("Groups of duplicates after shrink: %d\n", totalCount)

	if config.RemoveFlag {
		err := removeFiles(finalMap, logger, &config)
		if err != nil {
			return err
		}
	} else {
		logHashMap(logger, finalMap)
	}
	return nil

}

func logHashMap(logger *logging.Logger, hashMap map[string][]walkdir.FileInfo) {
	totalDupes := 0
	for hash, paths := range hashMap {
		dupeCount := len(paths)
		totalDupes += dupeCount
		logger.Info("duplicates_found",
			"hash", hash,
			"found", dupeCount,
			"files", paths,
		)
	}
	logger.Info("scan_complete", "total_duplicates", totalDupes)

}

func displayHashMap(hashMap map[string][]walkdir.FileInfo) {
	for hash, paths := range hashMap {
		dupeCount := len(paths)
		fmt.Printf("\nHash: %d, Duplicates: %d, Size: %d bytes\n", hash, dupeCount, paths[0].FileSize)
		for _, path := range paths {
			fmt.Printf(" - %s\n", path.FilePath)
		}
	}
}

func removeFiles(hashMap map[string][]walkdir.FileInfo, logger *logging.Logger, config *Config) error {
	// Setup input for user choices for delete, remove, etc
	tty, err := os.Open("/dev/tty")
	if err != nil {
		return fmt.Errorf("cannot open tty for interactive input: %v", err)
	}
	defer tty.Close()

	reader := bufio.NewReader(tty)

nextHash:
	for hash, paths := range hashMap {

		subMap := map[string][]walkdir.FileInfo{
			hash: paths,
		}
		//display list of files with this same hash
		displayHashMap(subMap)

		// iterate through file list
	nextDuplicate:
		for _, file := range paths {

			fmt.Printf("Remove file: %s?\n", file.FilePath)

			choice, err := getUserChoice(reader)
			if err != nil {
				logger.Error("choice_err",
					"err", err,
				)
			}

			switch choice {
			case "d": //delete file
				err := os.Remove(file.FilePath)
				if err != nil {
					logger.Error("elete_error",
						"file_path", file.FilePath,
						"err", err,
					)
				} else {
					logger.Info("deleted_file",
						"file_path", file.FilePath,
					)
				}
				continue nextDuplicate

			case "t": //trash file
				err := trashFile(file.FilePath, hash, config)
				if err != nil {
					logger.Error("trash_error",
						"file_path", file.FilePath,
						"err", err,
					)
				} else {
					logger.Info("trashed_file",
						"file_path", file.FilePath,
					)
				}

			case "s": //skip file
				continue nextDuplicate

			case "c": //continue to next hash
				continue nextHash

			default:
				return nil

			}
		}

	}
	return nil
}

func copyFile(src, dest string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	return err
}

func configTrash() (string, string, error) {
	//Get username for trash path
	usr, err := user.Current()
	if err != nil {
		return "", "", fmt.Errorf("unable to get current user: %v", err)
	}

	// Define the trash path based on the OS
	// currently, only Linux will use the freedesktop spec path, other OSes will use a "trash" directory in the current working directory for simplicity
	var trashPath, trashInfoDir string

	if runtime.GOOS == "linux" {
		trashPath = filepath.Join(usr.HomeDir, ".local/share/Trash/files/")
		trashInfoDir = filepath.Join(usr.HomeDir, ".local/share/Trash/info/")
		// Ensure the trash info directory exists
		if _, err := os.Stat(trashInfoDir); os.IsNotExist(err) {
			err := os.MkdirAll(trashInfoDir, 0755)
			if err != nil {
				return "", "", fmt.Errorf("error creating trash info directory: %v", err)
			}
		}
	} else {
		trashPath = "trash"
		trashInfoDir = "trash"
		// Ensure the trash directory exists
		if _, err := os.Stat(trashPath); os.IsNotExist(err) {
			err := os.Mkdir(trashPath, 0755)
			if err != nil {
				return "", "", err
			}
		}
	}
	return trashPath, trashInfoDir, nil

}

func trashFile(filePath string, hashVal string, config *Config) error {

	// Create a unique name for the file in the trash to avoid conflicts
	ext := filepath.Ext(filePath)
	name := strings.TrimSuffix(filepath.Base(filePath), ext)
	enumeratedName := fmt.Sprintf("%s_%s%s", name, hashVal[:8], ext)

	destPath := filepath.Join(config.TrashPath, enumeratedName)
	src := filePath

	// Move the file to the trash, adding trashPath to the file name
	// First try to rename (move) the file, which is more efficient.
	err := os.Rename(filePath, destPath)
	if err != nil {
		// Rename failed, try copy + delete method as a fallback
		err = copyFile(src, destPath)
		if err != nil {
			return err
		}
		err = os.Remove(src)
		if err != nil {
			return err
		}
	}

	// Create .trashinfo file (to FreeDesktop spec) if on Linux in appropriate directory, non-Linux will place .trashinfo files
	// in the same directory as the trashed files for simplicity

	infoPath := filepath.Join(config.TrashInfoPath, enumeratedName+".trashinfo")
	originalPath := filePath
	infoContent := fmt.Sprintf("[Trash Info]\nPath=%s\nDeletionDate=%s\n", url.PathEscape(originalPath), time.Now().Format("2006-01-02T15:04:05"))

	err = os.WriteFile(infoPath, []byte(infoContent), 0644)
	if err != nil {
		return err
	}

	return nil
}

func getUserChoice(reader *bufio.Reader) (string, error) {
	choices := map[string]bool{
		"d": true,
		"t": true,
		"s": true,
		"c": true,
	}

	for {
		fmt.Printf("(D)elete, (T)rash, (S)kip, (C)ontinue to next hash --> ")

		input, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}

		input = strings.TrimSpace(strings.ToLower(input))

		if choices[input] {
			return input, nil
		}
	}
}
func filterDuplicates(hashMap map[string][]walkdir.FileInfo) (map[string][]walkdir.FileInfo, int) {
	finalMap := make(map[string][]walkdir.FileInfo)
	count := 0
	for hash, paths := range hashMap {
		if len(paths) == 1 {
			continue //unique, ignore
		}
		count++
		finalMap[hash] = paths
	}
	return finalMap, count
}

func getSpinner(count int) string {
	frames := []string{"|", "/", "-", "\\"}
	return frames[count%len(frames)]
}

var firstPassHashes = make(map[string][]walkdir.FileInfo)
var mu sync.Mutex
var wg sync.WaitGroup

func processPartialFile(file string, filesize int64, logger *logging.Logger) {
	defer wg.Done()
	partialHash, err := hashfile.PartialHash(file)
	if err != nil {
		logger.Error("partial_hash_error",
			"file_path", file,
			"err", err,
		)
		return
	}
	var fileInfo walkdir.FileInfo
	fileInfo.FilePath = file
	fileInfo.FileSize = filesize
	mu.Lock()
	firstPassHashes[partialHash] = append(firstPassHashes[partialHash], fileInfo)
	mu.Unlock()
}

var finalDuplicates = make(map[string][]walkdir.FileInfo)

func processFullFile(file string, filesize int64, logger *logging.Logger) {
	defer wg.Done()
	fullHash, err := hashfile.FullHash(file)
	if err != nil {
		logger.Error("full_hash_err",
			"file_path", file,
			"err", err,
		)
		return
	}
	var fileInfo walkdir.FileInfo
	fileInfo.FilePath = file
	fileInfo.FileSize = filesize
	mu.Lock()
	finalDuplicates[fullHash] = append(finalDuplicates[fullHash], fileInfo)
	mu.Unlock()
}
