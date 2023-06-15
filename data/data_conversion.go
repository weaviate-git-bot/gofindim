package data

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"

	"gocv.io/x/gocv"
)

func FormatDescriptors(desc *gocv.Mat) []byte {
	// descByte := make([]byte, desc.Rows()*desc.Cols())
	cols := desc.Cols()
	rows := desc.Rows()
	fmt.Printf("Rows: %d, Cols: %d\n", rows, cols)
	binaryVector := make([]byte, rows*cols)

	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			binaryVector[(i*cols)+j] = byte(desc.GetIntAt(i, j))
		}
	}

	return binaryVector
	// for i := 0; i < desc.Rows(); i++ {
	// 	for j := 0; j < desc.Cols(); j++ {
	// 		// data := make([]byte, 256)
	// 		// data[0] = byte(desc.GetUCharAt(i, j))
	// 		descByte = append(descByte, desc.GetUCharAt(i, j))
	// 	}
	// }
	// return descByte
}

func MatToBinaryVector(mat *gocv.Mat) []byte {
	dst := gocv.NewMat()
	mat.ConvertTo(&dst, gocv.MatTypeCV8U)
	data := dst.ToBytes() // Flatten Mat into a slice of bytes
	return data
	// rows := desc.Rows()
	// cols := desc.Cols()
	// descByte := make([]byte, rows*cols)
	// for i := 0; i < rows; i++ {
	// 	for j := 0; j < cols; j++ {
	// 		descByte = append(descByte, desc.GetUCharAt(i, j))
	// 	}
	// }
	// return descByte
}

func KeyPointsToFloat32(kp *[]gocv.KeyPoint) []float32 {
	var kpFloat []float32
	for _, k := range *kp {
		kpFloat = append(kpFloat, float32(k.X))
		kpFloat = append(kpFloat, float32(k.Y))
		kpFloat = append(kpFloat, float32(k.Size))
		kpFloat = append(kpFloat, float32(k.Angle))
		kpFloat = append(kpFloat, float32(k.Response))
		kpFloat = append(kpFloat, float32(k.Octave))
		kpFloat = append(kpFloat, float32(k.ClassID))
	}
	return kpFloat
}

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
