package data

import (
	"encoding/base64"
	"image"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/corona10/goimagehash"
	"gocv.io/x/gocv"
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
	vector, err := i.ToVector()
	if err != nil {
		log.Println("Error converting vector during conversion to interface")
	}
	return map[string]interface{}{
		"name":      i.Name,
		"path":      i.Path,
		"embedding": vector,
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
	path, err = filepath.Abs(path)
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

func (i *ImageFile) ToVector() ([]float32, error) {
	vector, err := VectorizeImage(i)
	if err != nil {
		return nil, err
	}
	i.Embedding = vector
	return vector, nil
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
