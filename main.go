package main

import (
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
	Path    string
	Trash   bool
	Delete  bool
	Verbose bool
	LogPath string
}

func main() {
	// Check if data is being piped in through stdin
	/* not yet implimented
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		// Data is being piped in
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			path := scanner.Text()
			if path == "" {
				continue
			}

	}
			// Now you can process this path just like a command-line argument!
		}
			} else {
				// No pipe, maybe look for command line arguments?
			}
		}
	} */

	// Define flags
	trashFlag := flag.Bool("trash", false, "Trash duplicate files instead of just listing")
	deletFlag := flag.Bool("delete", false, "Delete duplicate files instead of just listing")
	logFlag := flag.String("log", "none", "Log path, or 'default' for current directory")
	verboseFlag := flag.Bool("v", false, "verbose mode,")

	flag.Parse()

	// Remaining arguments after flags are parsed
	args := flag.Args()
	if len(args) < 1 {
		log.Fatalf("Usage: %s [flags] <path>\n", os.Args[0])
	}
	targetPath := args[0]

	// Create config struct with parsed values
	config := Config{
		Path:    targetPath,
		Trash:   *trashFlag,
		Delete:  *deletFlag,
		Verbose: *verboseFlag,
		LogPath: *logFlag,
	}

	// Create logger. All output will be done through the logger, which will handle writing to file and/or screen based on config
	logger, err := logger.NewLogger(config.LogPath, config.Verbose)
	if err != nil {
		log.Fatalf("Error creating logger: %v", err)
	}
	defer logger.Close()

	err = process(config, logger)
	if err != nil {
		logger.Error("Error running directory mode: %v", err)
	}

}

func process(config Config, logger *logger.Logger) error {
	returnedMap, err := walkdir.WalkDir(config.Path, logger)
	if err != nil {
		return fmt.Errorf("Error walking directory: %v", err)
	}

	if config.Trash {
		if err := trashDuplicateFiles(returnedMap, logger); err != nil {
			return fmt.Errorf("Error trashing duplicate files: %v", err)
		}
	} else if config.Delete {
		if err := deleteDuplicateFiles(returnedMap, logger); err != nil {
			return fmt.Errorf("Error deleting duplicate files: %v", err)
		}
	} else {
		// just list duplicates, do nothing else
		displayHashMap(logger, returnedMap, config)
	}

	return nil
}

func displayHashMap(logger *logger.Logger, hashMap map[string][]walkdir.FileInfo, config Config) {
	for hash, paths := range hashMap {
		if config.Verbose {
			// if verbose, print all files, even if not duplicates, and include file sizes
			logger.Log("Hash: %s", hash)

			for _, path := range paths {
				logger.Log(" - %s size: %d", path.FilePath, path.FileSize)
			}

		} else { // if not verbose, just print instances with duplicates

			if len(paths) > 1 {
				logger.Log("Duplicate files with hash: %s", hash)
				for _, path := range paths {
					logger.Log(" - %s size: %d", path.FilePath, path.FileSize)

				}
			}
		}
	}
}

func trashDuplicateFiles(hashMap map[string][]walkdir.FileInfo, logger *logger.Logger) error {
	//Get username for trash path
	usr, err := user.Current()
	if err != nil {
		return fmt.Errorf("unable to get current user: %v", err)
	}

	// Define the trash path based on the OS
	var trashPath, trashInfoDir string

	if runtime.GOOS == "linux" {

		trashPath = filepath.Join(usr.HomeDir, usr.Username, ".local/share/Trash/files/")
		trashInfoDir = filepath.Join(usr.HomeDir, usr.Username, ".local/share/Trash/info/")
		// Ensure the trash info directory exists
		if _, err := os.Stat(trashInfoDir); os.IsNotExist(err) {
			err := os.MkdirAll(trashInfoDir, 0755)
			if err != nil {
				return fmt.Errorf("Error creating trash info directory: %v", err)
			}
		}
	} else {
		trashPath = "trash"
		trashInfoDir = "trash"
		// Ensure the trash directory exists
		if _, err := os.Stat(trashPath); os.IsNotExist(err) {
			err := os.Mkdir(trashPath, 0755)
			if err != nil {
				return err
			}
		}
	}

	for _, paths := range hashMap {
		if len(paths) > 1 {
			// Keep the first file and trash the rest
			// *Future improvement*: iterate through duplicates and ask user which one to keep,
			//  or if they want to keep all, trash all, etc. For now, just keep the first one and trash the rest.

			for i := 1; i < len(paths); i++ {

				// Create a unique name for the file in the trash to avoid conflicts
				ext := filepath.Ext(paths[i].FilePath)
				name := strings.TrimSuffix(filepath.Base(paths[i].FilePath), ext)
				enumeratedName := fmt.Sprintf("%s_%d%s", name, i, ext)

				destPath := filepath.Join(trashPath, enumeratedName)
				src := paths[i].FilePath

				// Move the file to the trash, adding trashPath to the file name
				// First try to rename (move) the file, which is more efficient.
				err := os.Rename(paths[i].FilePath, destPath)
				if err != nil {
					// Rename failed, try copy + delete method as a fallback
					err = copyFile(src, destPath)
					if err != nil {
						logger.Error("Error copying file to trash %s: %v", paths[i].FilePath, err)
						return err
					}
					err = os.Remove(src)
					if err != nil {
						logger.Error("Error deleting original file after copying to trash %s: %v", paths[i].FilePath, err)
						return err
					}

					logger.Log("Trashed file (copy+delete): %s", paths[i].FilePath)
				} else {
					logger.Log("Trashed file: %s", paths[i].FilePath)
				}

				// Create .trashinfo file (to FreeDesktop spec) if on Linux in appropriate directory, non-Linux will place .trashinfo files
				// in the same directory as the trashed files for simplicity

				infoPath := filepath.Join(trashInfoDir, enumeratedName+".trashinfo")
				originalPath := paths[i].FilePath
				infoContent := fmt.Sprintf("[Trash Info]\nPath=%s\nDeletionDate=%s\n", url.PathEscape(originalPath), time.Now().Format("2006-01-02T15:04:05"))

				err = os.WriteFile(infoPath, []byte(infoContent), 0644)
				if err != nil {
					logger.Error("Error creating trash info file for %s: %v", paths[i].FilePath, err)
					return err
				}

			}
		}
	}
	return nil
}

func deleteDuplicateFiles(hashMap map[string][]walkdir.FileInfo, logger *logger.Logger) error {
	// Iterate through the hash map and delete duplicate files, keeping the first instance
	// *Future improvement*: iterate through duplicates and ask user which one to keep.
	for _, paths := range hashMap {
		if len(paths) > 1 {
			// Keep the first file and delete the rest
			for i := 1; i < len(paths); i++ {
				err := os.Remove(paths[i].FilePath)
				if err != nil {
					logger.Error("Error deleting file %s: %v", paths[i].FilePath, err)
				} else {
					logger.Log("Deleted duplicate file: %s", paths[i].FilePath)
				}
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
