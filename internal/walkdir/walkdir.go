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
	"log"
	"os"
	"path/filepath"

	hashfile "github.com/gainax2k1/hash-file-compare/internal/hashfile"
)

type FileInfo struct {
	FilePath string
	FileSize int64
}

//takes in a directory path, returns a map of hash values (of each file),
// and slice of the file paths that correspond to each hash value

func WalkDir(dir string) (map[string][]FileInfo, error) {
	// map to store hash values and corresponding file paths
	hashMap := make(map[string][]FileInfo)

	//refactoring using WalkDir function (replaced deprecated Walk function)
	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// Process only files, ignore directories
		if !d.IsDir() {
			hashValue, err := hashfile.HashFromFilename(path)
			// fmt.Println(hashValue, path)
			if err != nil {
				log.Printf("Error hashing file %s: %v\n", path, err)
				return nil // Continue with next file
			}
			fileSize, err := getFileSize(path)
			if err != nil {
				log.Printf("Error getting file size for %s: %v\n", path, err)
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
		return nil, fmt.Errorf("Error walking: %w", err)
	}
	// Check if any files were processed
	if len(hashMap) == 0 {
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
