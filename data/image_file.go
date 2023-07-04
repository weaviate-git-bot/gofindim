package data

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/corona10/goimagehash"
)

type ImageFile struct {
	image.Image
	Format    string
	Name      string
	Path      string
	Base64    string
	Embedding []float32
}

func (m *ImageFile) ToWeaviateImageData() map[string]interface{} {
	imageDataInterface := map[string]interface{}{
		"filename": m.Name,
		"rating":   5, // sample rating
	}
	return imageDataInterface
}

func (i *ImageFile) toInterface() map[string]interface{} {
	var err error
	// vector, err := i.ToVector()
	if err != nil {
		log.Println("Error converting vector during conversion to interface")
	}
	return map[string]interface{}{
		"name": i.Name,
		"path": i.Path,
		// "embedding": vector,
	}

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
	// path, err = filepath.Abs(path)
	filename := filepath.Base(file.Name())
	_img, format, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	raw, err := ioutil.ReadFile(path)
	b64 := base64.StdEncoding.EncodeToString(raw)
	myImg := &ImageFile{Image: _img, Format: format, Name: filename, Path: path, Base64: b64}
	return myImg, nil
}

func NewImageFileFromFormFile(img multipart.File, name string) *ImageFile {
	_img, format, err := image.Decode(img)
	img.Seek(0, io.SeekStart)
	if err != nil {
		log.Fatal(err)
	}
	println("format", format)
	img.Seek(0, io.SeekStart)
	raw, err := io.ReadAll(img)
	if err != nil {
		log.Fatal(err)
	}
	println("About to encode")
	b64 := base64.StdEncoding.EncodeToString(raw)
	if len(b64) == 0 {
		panic("b64 is empty")
	}
	myImg := &ImageFile{Image: _img, Format: format, Name: name, Base64: b64}
	return myImg
}

func NewImageFileFromURL(url string, name string) (*ImageFile, error) {
	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error fetching image:", err)
		return nil, err
	}
	defer response.Body.Close()

	// Check if the request was successful (status code 200)
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status: %s", response.Status)
	}

	// Read the response data (image content)
	imageData, err := ioutil.ReadAll(response.Body)
	imageReader := bytes.NewReader(imageData)
	_img, format, err := image.Decode(imageReader)

	imageReader.Seek(0, io.SeekStart)
	raw, err := io.ReadAll(imageReader)
	if err != nil {
		return nil, err
	}
	println("About to encode")
	b64 := base64.StdEncoding.EncodeToString(raw)
	if len(b64) == 0 {
		panic("b64 is empty")
	}
	myImg := &ImageFile{Image: _img, Format: format, Name: name, Base64: b64}
	return myImg, nil
}

func (i *ImageFile) ToVector() ([]float32, error) {
	vector, err := VectorizeImage(i)
	if err != nil {
		return nil, err
	}
	i.Embedding = vector
	return vector, nil
}
