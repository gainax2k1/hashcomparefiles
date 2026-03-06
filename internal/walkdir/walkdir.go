package walkdir

import (
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
		return nil, fmt.Errorf("error resolving absolute path: %w", err)
	}

	// map to store hash values and corresponding file paths
	hashMap := make(map[string][]FileInfo)

	//refactored using WalkDir function (replaced deprecated Walk function)
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

			info, err := d.Info() // replaces my getFileSize(path) call, which was redundant (os.DirEntry already has it)
			if err != nil {
				logger.Error("Error getting file info for %s: %v", path, err)
				return nil // Skip this file and move to the next
			}

			fileSize := info.Size()
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

			// appends file to slice of files with that hash value
			hashMap[hashValue] = append(hashMap[hashValue], fileInfo)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("Error walking: %w", err)
	}

	// Return the map of hash values and file paths
	return hashMap, nil
}
