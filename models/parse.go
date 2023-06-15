package models

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/corona10/goimagehash"
	"gocv.io/x/gocv"
)

type ImageFile struct {
	image.Image
	Format string
	Name   string
	Path   string
}

func (m *ImageFile) ToWeaviateImageData() map[string]interface{} {
	imageDataInterface := map[string]interface{}{
		"filename": m.Name,
		"rating":   5, // sample rating
	}
	return imageDataInterface
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
func printMat(mat gocv.Mat) {
	rows := mat.Rows()
	cols := mat.Cols()
	var output strings.Builder

	for row := 0; row < rows; row++ {
		for col := 0; col < cols; col++ {
			value := mat.GetFloatAt(row, col)
			output.WriteString(fmt.Sprintf("%f ", value))
		}
		output.WriteString("\n")
	}

	fmt.Print(output.String())
}
func GetImageSiftDescriptors(img1 *ImageFile) ([]gocv.KeyPoint, gocv.Mat) {
	mat1 := img1.toImgMat()
	defer mat1.Close()

	sift := gocv.NewSIFT()
	defer sift.Close()

	keypoints, descriptors := sift.DetectAndCompute(*mat1, gocv.NewMat())
	return keypoints, descriptors
}

func CompareDescriptors(img1, img2 *ImageFile) (float64, error) {
	bfmatcher := gocv.NewBFMatcherWithParams(gocv.NormL2, false)
	defer bfmatcher.Close()
	_, mat1 := GetImageSiftDescriptors(img1)
	_, mat2 := GetImageSiftDescriptors(img2)

	matches := bfmatcher.KnnMatch(mat1, mat2, 2)
	goodMatches := 0
	for _, mPair := range matches {
		if len(mPair) == 2 && mPair[0].Distance < 0.75*mPair[1].Distance {
			goodMatches++
		}
	}
	return float64(goodMatches) / float64(len(matches)), nil
}

func CompareImageOrb(image1, image2 *ImageFile) (int, error) {
	// Initiate ORB detector
	orb := gocv.NewORBWithParams(128, 1.2, 8, 31, 0, 2, gocv.ORBScoreTypeHarris, 31, 20)

	defer orb.Close()

	_, desc1 := image1.toKeypointsDescriptors(&orb)
	_, desc2 := image2.toKeypointsDescriptors(&orb)
	// printMat(desc1)
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
