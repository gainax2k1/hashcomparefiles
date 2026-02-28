package main

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"runtime"

	hashfile "github.com/gainax2k1/hash-file-compare/hashFile"
	walkDir "github.com/gainax2k1/hash-file-compare/walkDir"
)

func main() {
	fmt.Println("Find duplicate files by hash value")
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s <filename>\n", os.Args[0])
	}

	// TODO!: Add logging, for duplicate files/deleteted/trashed files, and errors. Maybe add a -v flag for verbose logging?

	// TODO!: Organize better flag/CLI options implimentation!

	// check for -d flag here to call WalkDir
	if os.Args[1] == "-d" {
		// verify there's a directory path argument
		if len(os.Args) < 3 {
			log.Fatalf("Usage: %s -d <directory_path>\n", os.Args[0])
		}

		// call WalkDir with the provided directory path
		returnedMap, err := walkDir.WalkDir(os.Args[2])
		if err != nil {
			log.Fatalf("Error walking directory: %v\n", err)
		}

		// Print hash files for debugging purposes
		/*
			for hash := range returnedMap {
				for _, path := range returnedMap[hash] {
					fmt.Printf("Hash: %s\nFiles: %v\n", hash, path)
				}
			}*/

		// Display duplicate files for debugging purposes
		fmt.Println("Printing duplicate files:")
		displayDupicateFiles(returnedMap)
		return
	}

	if os.Args[1] == "-TRASH" {
		// todo: add "/info" folder in .trash for storing info about trashed files, like original path, deletion date, etc. maybe add a -v flag to print this info when trashing files?
		/*
			example: for file called "3DMark.2.trashinfo"
			[Trash Info]
			Path=/mnt/nvme1n1p1/WIndows%20user%20folders/Pictures/3DMark
			DeletionDate=2026-02-08T19:53:54


		*/

		// verify there's a directory path argument
		if len(os.Args) < 3 {
			log.Fatalf("Usage: %s -TRASH <directory_path>\n", os.Args[0])
		}

		// call WalkDir with the provided directory path
		returnedMap, err := walkDir.WalkDir(os.Args[2])
		if err != nil {
			log.Fatalf("Error walking directory: %v\n", err)
		}

		// Remove duplicate files
		err = trashDuplicateFiles(returnedMap)
		if err != nil {
			log.Fatalf("Error trashing duplicate files: %v\n", err)
		}

		return
	}

	if os.Args[1] == "-DELETE" {
		// verify there's a directory path argument
		if len(os.Args) < 3 {
			log.Fatalf("Usage: %s -DELETE <directory_path>\n", os.Args[0])
		}

		// call WalkDir with the provided directory path
		returnedMap, err := walkDir.WalkDir(os.Args[2])
		if err != nil {
			log.Fatalf("Error walking directory: %v\n", err)
		}

		// Remove duplicate files
		deleteDuplicateFiles(returnedMap)
		return
	}

	// handles single file hash value check
	filename := os.Args[1]

	fileHashValue, err := hashfile.HashFile(filename)
	if err != nil {
		log.Fatalf("Error hashing file: %v\n", err)
	}

	fmt.Println(fileHashValue)

}

func displayDupicateFiles(hashMap map[string][]walkDir.FileInfo) {
	for hash, paths := range hashMap {
		if len(paths) > 1 {
			fmt.Printf("Hash: %s", hash)
			fmt.Println("Files:")
			for _, path := range paths {
				fmt.Printf(" path: %s size: %d\n", path.FilePath, path.FileSize)

			}
		}
	}
}

/*
NEED TO DO: Add functionality for linux (at least? windows might not be a problem) for handling trashing files
on other drives (currently only works on root drives). Maybe copy them to root drive's trash? or maybe move to a folder on that drive, label it as trash
and let user handle it?
*/
func trashDuplicateFiles(hashMap map[string][]walkDir.FileInfo) error {
	//Get username for trash path
	usr, err := user.Current()
	if err != nil {
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
					log.Printf("Error moving to trash file %s: %v\n", paths[i].FilePath, err)
					return err
				} else {
					fmt.Printf("Trashed file: %s\n", paths[i].FilePath)
				}
			}
		}
	}
	return nil
}

func deleteDuplicateFiles(hashMap map[string][]walkDir.FileInfo) {
	for _, paths := range hashMap {
		if len(paths) > 1 {
			// Keep the first file and delete the rest
			for i := 1; i < len(paths); i++ {
				err := os.Remove(paths[i].FilePath)
				if err != nil {
					log.Printf("Error deleting file %s: %v\n", paths[i].FilePath, err)
				} else {
					fmt.Printf("Deleted duplicate file: %s\n", paths[i].FilePath)
				}
			}
		}
	}
}
