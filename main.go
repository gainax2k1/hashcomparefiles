package main

import (
	"flag" // for future use with CLI options
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

	hashfile "github.com/gainax2k1/hashcomparefiles/internal/hashfile"
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
	// Define flags
	oneFile := flag.Bool("f", false, "Hash single file mode")
	dirMode := flag.Bool("d", false, "Hash directory mode")
	trashFlag := flag.Bool("trash", false, "Trash duplicate files instead of just listing")
	deletFlag := flag.Bool("delete", false, "Delete duplicate files instead of just listing")
	logFlag := flag.String("log", "none", "Log path, or 'default' for current directory")
	verboseFlag := flag.Bool("v", false, "verbose mode,")
	//list mode will be default behavior, just list duplicates without deleting or trashing.
	//listFlag := flag.Bool("list", false, "List duplicates without deleting/trashing")

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

	switch {
	case *oneFile:

		if err := runSingleFileMode(config, logger); err != nil {
			log.Fatal(err)
		}
	case *dirMode:
		if err := runDirectoryMode(config, logger); err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatal("No valid mode specified. Use -f for single file or -d for directory mode.")
	}

}

func runSingleFileMode(config Config, logger *logger.Logger) error { // Everywhere else, just use it:
	//logger.Log("Scanning: %s", config.Path)
	//logger.Log("Found duplicate: %s", filePath)
	//logger.Log("Error: %v", err)

	// for single files, ignore trash, delete, and verbose flags, just print the hash value and optionally log it
	fileHashValue, err := hashfile.HashFromFilename(config.Path)
	if err != nil {
		logger.Error("Error hashing file: %v", err)
		return fmt.Errorf("Error hashing file: %v", err)
	}
	//"hash-file-compare_%s.log", time.Now().Format("2006-01-02_150405"
	logger.Log("File: %s, Hash: %s", config.Path, fileHashValue)
	return nil
}

func runDirectoryMode(config Config, logger *logger.Logger) error {
	returnedMap, err := walkdir.WalkDir(config.Path, logger)
	if err != nil {
		logger.Error("Error walking directory: %v", err)
		return fmt.Errorf("Error walking directory: %v", err)
	}

	if config.Trash {
		if err := trashDuplicateFiles(returnedMap, logger); err != nil {
			logger.Error("Error trashing duplicate files: %v", err)
			return fmt.Errorf("Error trashing duplicate files: %v", err)
		}
	} else if config.Delete {
		if err := deleteDuplicateFiles(returnedMap, logger); err != nil {
			logger.Error("Error deleting duplicate files: %v", err)
			return fmt.Errorf("Error deleting duplicate files: %v", err)
		}
	} else {
		// just list duplicates, do nothing else
		displayDupicateFiles(logger, returnedMap, config)
	}

	return nil
}

func displayDupicateFiles(logger *logger.Logger, hashMap map[string][]walkdir.FileInfo, config Config) {
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
		logger.Error("unable to get current user: %v", err)
		return fmt.Errorf("unable to get current user: %v", err)
	}

	// Define the trash path based on the OS
	var trashPath, trashInfoDir string

	if runtime.GOOS == "linux" {

		trashPath = filepath.Join("/home", usr.Username, ".local/share/Trash/files/")
		trashInfoDir = filepath.Join("/home", usr.Username, ".local/share/Trash/info/")
		// Ensure the trash info directory exists
		if _, err := os.Stat(trashInfoDir); os.IsNotExist(err) {
			err := os.MkdirAll(trashInfoDir, 0755)
			if err != nil {
				logger.Error("Error creating trash info directory: %v", err)
				return fmt.Errorf("Error creating trash info directory: %v", err)
			}
		}
	} else {
		trashPath = "trash"
		trashInfoDir = ""
		// Ensure the trash directory exists
		if _, err := os.Stat(trashPath); os.IsNotExist(err) {
			err := os.Mkdir(trashPath, 0755)
			if err != nil {
				logger.Error("Error creating trash info directory: %v", err)
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
				err := os.Rename(paths[i].FilePath, destPath)
				if err != nil {
					// Rename failed, try copy + delete method as a fallback (e.g. if moving across different filesystems)
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

				// Create .trashinfo file
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
