package data

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
)

// Hash a file to get its UUID
func FileToUUID(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hasher := sha256.New()
	_, err = io.Copy(hasher, file)
	if err != nil {
		return "", err
	}

	hash := hasher.Sum(nil)
	return hex.EncodeToString(hash), nil
}
