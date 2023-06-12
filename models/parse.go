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
	"gocv.io/x/gocv"
)

type ImageFile struct {
	image.Image
	Format string
	Name   string
	Path   string
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
	path, err = filepath.Abs(path)
	filename := filepath.Base(file.Name())
	_img, format, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	myImg := &ImageFile{Image: _img, Format: format, Name: filename, Path: path}
	return myImg, nil
}

func (i *ImageFile) toImgMat() *gocv.Mat {
	imgMat := gocv.IMRead(i.Path, gocv.IMReadColor)
	return &imgMat
}

func (i *ImageFile) toMat() *gocv.Mat {
	imgMat := i.toImgMat()
	gray := gocv.NewMat()
	defer gray.Close()
	gocv.CvtColor(*imgMat, &gray, gocv.ColorBGRToGray)
	return imgMat
}

func (i *ImageFile) toKeypointsDescriptors(orb *gocv.ORB) ([]gocv.KeyPoint, gocv.Mat) {
	imgMat := i.toMat()
	defer imgMat.Close()
	kp, desc := orb.DetectAndCompute(*imgMat, gocv.NewMat())
	return kp, desc
}

func CompareImageOrb(image1, image2 *ImageFile) (int, error) {
	// Initiate ORB detector
	orb := gocv.NewORB()
	defer orb.Close()

	_, desc1 := image1.toKeypointsDescriptors(&orb)
	_, desc2 := image2.toKeypointsDescriptors(&orb)
	matcher := gocv.NewBFMatcherWithParams(gocv.NormHamming, false)
	matches := matcher.KnnMatch(desc1, desc2, 2)
	goodMatches := make([]gocv.DMatch, 0)
	for _, m := range matches {
		if len(m) < 2 {
			continue
		}
		if m[0].Distance < 0.75*m[1].Distance {
			goodMatches = append(goodMatches, m[0])
		}
	}
	return len(goodMatches), nil
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
