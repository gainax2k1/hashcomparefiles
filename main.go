package main

import (
	"flag" // for future use with CLI options
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"runtime"

	hashfile "github.com/gainax2k1/hash-file-compare/internal/hashfile"
	"github.com/gainax2k1/hash-file-compare/internal/logger"
	walkdir "github.com/gainax2k1/hash-file-compare/internal/walkdir"
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

/*
NEED TO DO: Add functionality for linux (at least? windows might not be a problem) for handling trashing files
on other drives (currently only works on root drives). Maybe copy them to root drive's trash? or maybe move to a folder on that drive, label it as trash
and let user handle it?
*/
// todo: add "/info" folder in .trash for storing info about trashed files, like original path, deletion date, etc. maybe add a -v flag to print this info when trashing files?
/*
	example: for file called "3DMark.2.trashinfo"
	[Trash Info]
	Path=/mnt/nvme1n1p1/WIndows%20user%20folders/Pictures/3DMark
	DeletionDate=2026-02-08T19:53:54


*/

func trashDuplicateFiles(hashMap map[string][]walkdir.FileInfo, logger *logger.Logger) error {
	//Get username for trash path
	usr, err := user.Current()
	if err != nil {
		logger.Error("unable to get current user: %v", err)
		return fmt.Errorf("unable to get current user: %v", err)
	}

	// Define the trash path based on the OS
	var trashPath string
	switch runtime.GOOS {
	case "windows":
		trashPath = "C:\\$Recycle.Bin\\"
	case "darwin": //macOS
		trashPath = filepath.Join("/Users", usr.Username, ".Trash")
	case "linux":
		trashPath = filepath.Join("/home", usr.Username, ".local/share/Trash/files/")
	default:
		return fmt.Errorf("unsupported OS for trashing files")
	}

	for _, paths := range hashMap {
		if len(paths) > 1 {
			// Keep the first file and trash the rest
			for i := 1; i < len(paths); i++ {

				destPath := filepath.Join(trashPath, filepath.Base(paths[i].FilePath))
				// Move the file to the trash, adding trashPath to the file name
				err := os.Rename(paths[i].FilePath, destPath)

				if err != nil {

					logger.Error("Error moving to trash file %s: %v", paths[i].FilePath, err)
					return err
				} else {
					logger.Log("Trashed file: %s", paths[i].FilePath)
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
