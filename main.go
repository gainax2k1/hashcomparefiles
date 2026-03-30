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
	"time"

	"github.com/gainax2k1/hashcomparefiles/internal/hashfile"
	"github.com/gainax2k1/hashcomparefiles/internal/logger"
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
	logFlag := flag.String("log", "none", "Log filename, or 'default' for current directory log.log")
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
	logger, err := logger.NewLogger(config.LogPath, *verboseFlag)
	if err != nil {
		log.Fatalf("Error creating logger: %v", err)
	}

	defer logger.Close()

	logger.Log("(Start)")

	err = process(targets, config, logger)
	if err != nil {
		logger.Error("Error processing: %v", err)
	}
	logger.Log("(Done)")

}

func process(targets []string, config Config, logger *logger.Logger) error {
	// 1. Make map (key=filesize, value=[]filepaths)
	// --- ignore all symlinks, zero size files
	// 2. For each key, if len(value) > 1, then run smaller hash on each file, make map of (key=hash, value=[]filepaths)
	// --- first pass, save on hashing on large files unless neccessary
	// 3. For each key, if len(value) > 1, then run full hash on each file, make map of (key=hash, value=[]filepaths)
	// --- 2nd pass, run fullhash on remaining files

	var spinnerCounter = 0

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
			if size < config.Maxflag && size > config.Minflag {
				fileSizeMap[size] = append(fileSizeMap[size], files...)
			}
		}
	}
	logger.Log("Filecount after pass (1/3): %d", totalCount)
	fmt.Printf("Filecount after pass (1/3): %d\n", totalCount)

	// SECOND PASS:
	firstPassHashes := make(map[string][]walkdir.FileInfo)
	totalCount = 0 //reset

	for filesize, files := range fileSizeMap {
		if len(files) == 1 { // only one file with this size, so unique
			continue // skip this file
		}
		// multiple files with this size, so we need to compare them
		for _, file := range files {

			spinnerCounter++
			if spinnerCounter%100 == 0 {
				// Pulse/Spinner every 100 files to save CPU
				// \r clears the line, then we print the spinner and count
				fmt.Fprintf(os.Stderr, "\r %s Files processed: %d\r", getSpinner(spinnerCounter/100), spinnerCounter)

			}

			partialHash, err := hashfile.PartialHash(file)
			if err != nil {
				logger.Error("Error partial hashing file %s: %v", file, err)
				continue // skip this file
			}
			var fileInfo walkdir.FileInfo
			fileInfo.FilePath = file
			fileInfo.FileSize = filesize

			firstPassHashes[partialHash] = append(firstPassHashes[partialHash], fileInfo)
			totalCount++
		}
	}
	logger.Log("Filecount after pass (2/3): %d", totalCount)

	fmt.Printf("Filecount after pass (2/3): %d\n", totalCount)
	spinnerCounter = 0

	// THIRD PASS:
	finalDuplicates := make(map[string][]walkdir.FileInfo)
	totalCount = 0 //reset
	for smallHash, files := range firstPassHashes {
		if len(files) == 1 { // only one file with this size
			continue // skip this file, hash to be unique
		}
		for _, file := range files {

			spinnerCounter++
			if spinnerCounter%100 == 0 {
				// Pulse/Spinner every 100 files to save CPU
				// \r clears the line, then we print the spinner and count
				fmt.Fprintf(os.Stderr, "\r %s Files processed: %d\r", getSpinner(spinnerCounter/100), spinnerCounter)
			}

			if file.FileSize <= hashfile.PARTIALBYTELENGTH {

				finalDuplicates[smallHash] = append(finalDuplicates[smallHash], file)
				totalCount++
				continue // use first hash, since file was already *fully* hashed
			}
			fullHash, err := hashfile.FullHash(file.FilePath)
			if err != nil {
				logger.Error("Error full hashing file %s: %v", file, err)
				continue // skip this file
			}
			finalDuplicates[fullHash] = append(finalDuplicates[fullHash], file)
			totalCount++
		}
	}

	logger.Log("Filecount after pass (3/3): %d", totalCount)

	fmt.Printf("Filecount after pass (3/3): %d\n", totalCount)

	//shrink map
	finalMap, totalCount := filterDuplicates(finalDuplicates)

	logger.Log("Groups of duplicates after shrink: %d", totalCount)
	fmt.Printf("Groups of duplicates after shrink: %d\n", totalCount)

	if config.RemoveFlag {
		err := removeFiles(finalMap, logger, &config)
		if err != nil {
			return err
		}
	} else {
		displayHashMap(logger, finalMap)
	}
	return nil

}

func displayHashMap(logger *logger.Logger, hashMap map[string][]walkdir.FileInfo) {
	for hash, paths := range hashMap {
		count := 0
		logger.Log("Files with hash: %s", hash)
		for _, path := range paths {
			count++
			logger.Log(" - %s size: %d", path.FilePath, path.FileSize)
		}
		logger.Log(" -- Duplicates: %d", count)
	}

}

func removeFiles(hashMap map[string][]walkdir.FileInfo, logger *logger.Logger, config *Config) error {
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
		displayHashMap(logger, subMap)

		// iterate through file list
	nextDuplicate:
		for _, file := range paths {

			fmt.Printf("Remove file: %s?\n", file.FilePath)

			choice, err := getUserChoice(reader)
			if err != nil {
				logger.Error("Error geting user choice: %v", err)
			}

			switch choice {
			case "d": //delete file
				err := os.Remove(file.FilePath)
				if err != nil {
					logger.Error("Error deleting file %s: %v", file.FilePath, err)
				} else {
					logger.Log("Deleted duplicate file: %s", file.FilePath)
				}
				continue nextDuplicate

			case "t": //trash file
				err := trashFile(file.FilePath, hash, config)
				if err != nil {
					logger.Error("Error deleting file %s: %v", file.FilePath, err)
				} else {
					logger.Log("Deleted duplicate file: %s", file.FilePath)
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
		fmt.Printf(" - (D)elete, (T)rash, (S)kip, (C)ontinue to next hash > ")

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
