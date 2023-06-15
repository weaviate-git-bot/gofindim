package ai

import (
	"io/ioutil"
	"log"

	"github.com/owulveryck/onnx-go"
	"github.com/owulveryck/onnx-go/backend/x/gorgonnx"
)

type InputKind int

const (
	Text InputKind = iota
	Image
)

func MakeModel(modelPath string, inputKind InputKind, inputText string, inputImage []float32) {
	backend := gorgonnx.NewGraph()
	model := onnx.NewModel(backend)
	b, _ := ioutil.ReadFile("../ai_assets/model.onnx")
	err := model.UnmarshalBinary(b)
	if err != nil {
		log.Fatal(err)
	}

	// Prepare input
	// var input *tensor.Dense

	// switch inputKind {
	// case Text:
	// 	input = prepareTextInput(inputText)
	// case Image:
	// 	input, err = prepareImageInput(inputImage)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// default:
	// 	return nil, fmt.Errorf("unknown inputKind: %v", inputKind)
	// }

	// model.SetInput(0, tensor.New(tensor.WithShape(1, 1, 28, 28), tensor.Of(tensor.Float32)))
	// err = backend.Run()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// output, err := model.GetOutputTensors()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(output[0])
}
