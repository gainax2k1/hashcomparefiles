package hashfile

/*
	goal: to recieve a filename and return the sha256 value (hex) from it.
	plan: use the crypto/sha256 package to compute the hash value of the file.
	steps:
	1. open the file using os.Open
	2. create a new sha256 hasher using sha256.New()
	3. copy the file contents to the hasher using io.Copy
	4. compute the hash value and return it as a hex string using hex.EncodeToString
*/

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
)

func HashFromFilename(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return "", err
	}

	hash := hasher.Sum(nil)

	return hex.EncodeToString(hash), nil
}
