package walkdir

// purpose: walk through all the files in a directory
// then call the hashFile mod with each file
// then each returned hash should be stored in a map
// and then returned.

// todo: add filesize info?
// todo: make struct to hold values instead of map[string][]string (might be cleaner?, can inlcude filesize info in struct too)

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	logger "github.com/gainax2k1/hashcomparefiles/internal/logger"

	hashfile "github.com/gainax2k1/hashcomparefiles/internal/hashfile"
)

type FileInfo struct {
	FilePath string
	FileSize int64
}

//takes in a directory path, returns a map of hash values (of each file),
// and slice of the file paths that correspond to each hash value
// Key is the hash value, value is a slice of file paths that have that hash value

func WalkDir(dir string, logger *logger.Logger) (map[string][]FileInfo, error) {
	// resolve absolute path of the directory
	absDir, err := filepath.Abs(dir)
	if err != nil {
		logger.Error("Error resolving absolute path for directory %s: %v", dir, err)
		return nil, fmt.Errorf("error resolving absolute path: %w", err)
	}

	// map to store hash values and corresponding file paths
	hashMap := make(map[string][]FileInfo)

	//refactoring using WalkDir function (replaced deprecated Walk function)
	err = filepath.WalkDir(absDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// Process only files, ignore directories, symlinks, and empty files
		if !d.IsDir() {
			if d.Type()&os.ModeSymlink != 0 {
				logger.Log("Skipping symlink: %s", path)
				return nil
			}

			fileSize, err := getFileSize(path)
			if err != nil {
				logger.Error("Error getting file size for %s: %v", path, err)
				return nil // Continue with next file
			}

			if fileSize == 0 {
				logger.Log("Skipping empty file: %s", path)
				return nil
			}

			hashValue, err := hashfile.HashFromFilename(path)
			if err != nil {
				logger.Error("Error hashing file %s: %v", path, err)
				return nil // Continue with next file
			}

			fileInfo := FileInfo{
				FilePath: path,
				FileSize: fileSize,
			}

			hashMap[hashValue] = append(hashMap[hashValue], fileInfo)
		}

		return nil
	})

	if err != nil {
		logger.Error("Error walking directory: %v", err)
		return nil, fmt.Errorf("Error walking: %w", err)
	}
	// Check if any files were processed
	if len(hashMap) == 0 {
		logger.Log("No files found in the directory: %s", dir)
		return nil, errors.New("no files found in the directory")
	}

	// Return the map of hash values and file paths
	return hashMap, nil
}

func getFileSize(filename string) (int64, error) {
	fi, err := os.Stat(filename)
	if err != nil {
		return 0, err
	}
	return fi.Size(), nil
}
