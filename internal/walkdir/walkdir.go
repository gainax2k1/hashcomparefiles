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

// depricated
func WalkDir(dir string, logger *logger.Logger, runHash bool) (map[string][]FileInfo, int, error) {
	// resolve absolute path of the directory
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return nil, 0, fmt.Errorf("error resolving absolute path: %w", err)
	}

	// map to store hash values and corresponding file paths
	hashMap := make(map[string][]FileInfo)
	var count int

	defer func() {
		//ensures newline if spinner ran to not disrupt display
		if count >= 100 {
			fmt.Fprintf(os.Stderr, "\n")
		}
	}()

	//refactored using WalkDir function (replaced deprecated Walk function)
	err = filepath.WalkDir(absDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// Process only files, ignore directories, symlinks, and empty files
		if !d.IsDir() {
			if d.Type()&os.ModeSymlink != 0 {
				//skip symlinks
				return nil
			}

			info, err := d.Info() // replaces my getFileSize(path) call, which was redundant (os.DirEntry already has it)
			if err != nil {
				logger.Error("Error getting file info for %s: %v", path, err)
				return nil // Skip this file and move to the next
			}

			fileSize := info.Size()
			if fileSize == 0 {
				//skip empty files
				return nil
			}

			count++

			if runHash {
				hashValue, err := hashfile.FullHash(path)
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
			if count%100 == 0 {
				// Pulse/Spinner every 100 files to save CPU
				// \r clears the line, then we print the spinner and count
				fmt.Fprintf(os.Stderr, "\r %s Files processed: %d", getSpinner(count/100), count)
			}

		}

		return nil
	})

	if err != nil {
		return nil, 0, fmt.Errorf("Error walking: %w", err)
	}

	// Return the map of hash values and file paths
	return hashMap, count, nil
}

func WalkGetFileSizes(dir string, logger *logger.Logger) (map[int64][]string, int, error) {
	// resolve absolute path of the directory
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return nil, 0, fmt.Errorf("error resolving absolute path: %w", err)
	}

	var count int
	//key:filesize, value:[]filepath
	fileSizeMap := make(map[int64][]string)
	defer func() {
		//ensures newline if spinner ran to not disrupt display
		if count >= 100 {
			fmt.Fprintf(os.Stderr, "\n")
		}
	}()

	err = filepath.WalkDir(absDir, func(path string, d os.DirEntry, err error) error {
		//err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			logger.Error("Error walking directory: %s", err)
			return err
		}
		if !d.IsDir() {
			if d.Type()&os.ModeSymlink != 0 {
				//skip symlinks
				return nil
			}

			info, err := d.Info() // replaces my getFileSize(path) call, which was redundant (os.DirEntry already has it)
			if err != nil {
				logger.Error("Error getting file info for %s: %v", path, err)
				return nil // Skip this file and move to the next
			}

			fileSize := info.Size()
			if fileSize == 0 {
				//skip empty files
				return nil
			}
			fileSizeMap[fileSize] = append(fileSizeMap[fileSize], path)
			count++

			if count%100 == 0 {
				// Pulse/Spinner every 100 files to save CPU
				// \r clears the line, then we print the spinner and count
				fmt.Fprintf(os.Stderr, "\r %s Files processed: %d", getSpinner(count/100), count)
			}
		}
		return nil
	})
	if err != nil {
		return nil, 0, fmt.Errorf("Error walking: %w", err)
	}

	return fileSizeMap, count, nil
}

func getSpinner(count int) string {
	frames := []string{"|", "/", "-", "\\"}
	return frames[count%len(frames)]
}
