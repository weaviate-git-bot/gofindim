package models

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
	"path/filepath"

	"github.com/corona10/goimagehash"
)

type ImageFile struct {
	image.Image
	Format string
	Name   string
}

func (m *ImageFile) PerceptionHash() (string, error) {
	hash, err := goimagehash.PerceptionHash(m.Image)
	if err != nil {
		return "", err
	}
	return hash.ToString(), nil
}

func (m *ImageFile) ToImageModel() *ImageModel {
	hash, err := m.PerceptionHash()
	if err != nil {
		log.Fatal(err)
	}
	return &ImageModel{
		Filename: m.Name,
		Hash:     hash,
	}
}

func NewImageFileFromPath(path string) (*ImageFile, error) {
	// Open the file
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	filename := filepath.Base(file.Name())
	_img, format, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	myImg := &ImageFile{Image: _img, Format: format, Name: filename}
	return myImg, nil
}

func HashImagePerception(path string) (string, error) {
	// Open the file
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Decode the image
	img, _, err := image.Decode(file)
	if err != nil {
		fmt.Println("Error decoding image")
		return "", err
	}

	// Calculate the hash
	hash, err := goimagehash.PerceptionHash(img)
	if err != nil {
		return "", err
	}

	return hash.ToString(), nil
}
