package ai

import (
	"fmt"
	"image"
	"image/jpeg"
	"io/ioutil"
	"log"
	"os"

	"github.com/disintegration/imaging"
	"github.com/owulveryck/onnx-go"
	"github.com/owulveryck/onnx-go/backend/x/gorgonnx"
	"gorgonia.org/tensor"
)

type InputKind int

const (
	Text InputKind = iota
	Image
)

// GetFloat32Image Returns []float32 representation of image file
func GetFloat32Image(fname string) ([]float32, error) {
	file, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, err := jpeg.Decode(file)
	if err != nil {
		return nil, err
	}

	return Image2Float32(img)
}

// Image2Float32 Returns []float32 representation of image.Image
func Image2Float32(img image.Image) ([]float32, error) {
	width := img.Bounds().Dx()
	height := img.Bounds().Dy()
	imgwh := width * height
	imgSize := imgwh * 3

	ans := make([]float32, imgSize)
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			r, g, b, _ := img.At(y, x).RGBA()
			rpix, gpix, bpix := float32(r>>8)/float32(255.0), float32(g>>8)/float32(255.0), float32(b>>8)/float32(255.0)
			ans[y+x*height] = rpix
			ans[y+x*height+imgwh] = gpix
			ans[y+x*height+imgwh+imgwh] = bpix
		}
	}
	return ans, nil
}
func resizeImage(img image.Image, width, height int) *image.NRGBA {
	// Resize the image using the imaging package
	resized := imaging.Resize(img, width, height, imaging.Lanczos)

	return resized
}
func prepareImageInput(img image.Image) *tensor.Dense {
	// Resize and normalize the image according to the model requirements
	resized := resizeImage(img, 224, 224) // Assuming the input size for the model is 224

	// Convert the image to a tensor
	floats, err := Image2Float32(resized)
	if err != nil {
		log.Fatal(err)
	}
	// return tensor.New(tensor.WithShape(1, 3, 224, 224), tensor.Of(tensor.Float32), tensor.WithBacking(floats))
	return tensor.New(tensor.WithShape(1, 3, 224, 224), tensor.WithBacking(floats))
}

func imgToFloats(img image.Image) []float32 {
	b := img.Bounds()
	width, height := b.Max.X, b.Max.Y

	// Assumes the image is in the format: RGB
	data := make([]float32, 3*width*height)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			data[y*width+x] = float32(r>>8) / 255.0
			data[width*height+y*width+x] = float32(g>>8) / 255.0
			data[2*width*height+y*width+x] = float32(b>>8) / 255.0
		}
	}

	return data
}
func uint32ToFloat32Slice(input []uint32) []float32 {
	output := make([]float32, len(input))
	for i, v := range input {
		output[i] = float32(v)
	}
	return output
}

func prepareTextInput(input string) (*tensor.Dense, error) {
	tokenized, err := Tokenize(input)
	if err != nil {
		return nil, fmt.Errorf("error while tokenizing input: %w", err)
	}
	tokenTensor := tensor.New(
		tensor.WithShape(1, len(tokenized)),
		tensor.Of(tensor.Int32),
		tensor.WithBacking(uint32ToFloat32Slice(tokenized)))
	return tokenTensor, nil
}

func MakeImageEmbedding(img image.Image) ([]float32, error) {
	// imgTensor := prepareImageInput(img)
	backend := gorgonnx.NewGraph()
	model := onnx.NewModel(backend)
	b, _ := ioutil.ReadFile("ai_assets/model.onnx")
	err := model.UnmarshalBinary(b)
	if err != nil {
		return nil, fmt.Errorf("error while unmarshalling model: %w", err)
	}
	dummyData := make([]float32, 1*3*224*224)
	for i := range dummyData {
		dummyData[i] = 0.5 // Assign any value, e.g., 0.5
	}
	dummyTensor := tensor.New(tensor.WithShape(1, 3, 224, 224), tensor.WithBacking(dummyData))

	err = model.SetInput(0, dummyTensor)
	if err != nil {
		return nil, fmt.Errorf("error while setting input: %w", err)

	}
	fmt.Printf("input: %v\n", dummyTensor.Shape())
	err = backend.Run()
	if err != nil {
		return nil, fmt.Errorf("error while running model: %w", err)
	}
	output, err := model.GetOutputTensors()
	if err != nil {
		return nil, fmt.Errorf("error while getting output: %w", err)
	}
	return output[0].Data().([]float32), nil
}

func MakeModel(modelPath string, inputKind InputKind, inputText string, inputImage []float32) {
	backend := gorgonnx.NewGraph()
	model := onnx.NewModel(backend)
	b, _ := ioutil.ReadFile("../ai_assets/model.onnx")
	err := model.UnmarshalBinary(b)
	if err != nil {
		log.Fatal(err)
	}

}
