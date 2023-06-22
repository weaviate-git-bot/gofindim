package data

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
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
	// Get the hash digest and truncate or modify it to 32-length string
	hash := hasher.Sum(nil)[:16]
	unformattedUUID := hex.EncodeToString(hash)

	// Insert dashes to get the standard UUID format (8-4-4-12)
	formattedUUID := fmt.Sprintf("%s-%s-%s-%s-%s",
		unformattedUUID[0:8],
		unformattedUUID[8:12],
		unformattedUUID[12:16],
		unformattedUUID[16:20],
		unformattedUUID[20:32])

	return formattedUUID, nil
}
