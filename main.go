package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/gainax2k1/hashcomparefiles/internal/logger"
	walkdir "github.com/gainax2k1/hashcomparefiles/internal/walkdir"
)

type Config struct {
	FilePath         string
	TrashPath        string
	TrashInfoPath    string
	LogPath          string
	RemoveFlag       bool
	ShowPreHashCount bool
}

func main() {
	// Define flags and parse
	removeFlag := flag.Bool("remove", false, "Selectively choose which duplicates to trash or delete if desired")
	logFlag := flag.String("log", "none", "Log path, or 'default' for current directory")
	showPreHashCountFlag := flag.Bool("p", false, "Show Pre-hash file count (Potentially usefull for large runs, but now hits storage twice)")

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
		TrashPath:        trashPath,
		TrashInfoPath:    trashInfoPath,
		LogPath:          *logFlag,
		ShowPreHashCount: *showPreHashCountFlag,
		RemoveFlag:       *removeFlag,
	}

	// All output will be done through the logger, writing to file and/or screen based on config
	logger, err := logger.NewLogger(config.LogPath)
	if err != nil {
		log.Fatalf("Error creating logger: %v", err)
	}

	defer logger.Close()

	err = process(targets, config, logger)
	if err != nil {
		logger.Error("Error processing: %v", err)
	}
	logger.Log("(Done)")

}

func process(targets []string, config Config, logger *logger.Logger) error {

	masterMap := make(map[string][]walkdir.FileInfo)
	totalCount := 0
	runHash := false

	if config.ShowPreHashCount {
		for _, path := range targets {
			// Run the walk for each path and merge results into masterMap
			_, count, err := walkdir.WalkDir(path, logger, runHash)
			if err != nil {
				logger.Error("Skipping %s due to error: %v", path, err)
				continue // Keep going with other targets!
			}

			totalCount += count
		}
		logger.Log("Total files to process: %d", totalCount)
		totalCount = 0 // reset for hashing run
	}

	runHash = true
	for _, path := range targets {
		// Run the walk for each path and merge results into masterMap
		dirMap, count, err := walkdir.WalkDir(path, logger, runHash)
		if err != nil {
			logger.Error("Skipping %s due to error: %v", path, err)
			continue // Keep going with other targets!
		}

		totalCount += count
		// Merge dirMap into masterMap
		for hash, files := range dirMap {
			masterMap[hash] = append(masterMap[hash], files...)
		}
	}

	//shrink map to only duplicates because we don't need unique hashes
	masterMap, totalCount = filterDuplicates(masterMap)

	if config.RemoveFlag {
		err := removeFiles(masterMap, logger, &config)
		if err != nil {
			return err
		}
	} else {
		displayHashMap(logger, masterMap)
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

		//get counts for this hash
		pathsCount := len(paths)

		subMap := map[string][]walkdir.FileInfo{
			hash: paths,
		}
		//display list of files with this same hash
		displayHashMap(logger, subMap)

		// iterate through file list
	nextDuplicate:
		for i := 0; i < pathsCount; i++ {

			fmt.Printf("Delete file: %s?\n", paths[i].FilePath)

			choice, err := getUserChoice(reader)
			if err != nil {
				logger.Error("Error geting user choice: %v", err)
			}

			switch choice {
			case "d": //delete file
				err := os.Remove(paths[i].FilePath)
				if err != nil {
					logger.Error("Error deleting file %s: %v", paths[i].FilePath, err)
				} else {
					logger.Log("Deleted duplicate file: %s", paths[i].FilePath)
				}
				continue nextDuplicate

			case "t": //trash file
				err := trashFile(paths[i].FilePath, hash, config)
				if err != nil {
					logger.Error("Error deleting file %s: %v", paths[i].FilePath, err)
				} else {
					logger.Log("Deleted duplicate file: %s", paths[i].FilePath)
				}

			case "s": //skip file
				continue nextDuplicate

			case "c": //continue to next hash
				continue nextHash

			default:
				return nil //? shouldn't reach this...?

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
				return "", "", fmt.Errorf("Error creating trash info directory: %v", err)
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
	//enumeratedName := fmt.Sprintf("%s_%d%s", name, i, ext)
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
	dupesMap := make(map[string][]walkdir.FileInfo)
	count := 0
	for hash, paths := range hashMap {
		if len(paths) > 1 {

			dupesMap[hash] = paths
			count++
		}
	}
	return dupesMap, count
}
