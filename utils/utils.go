package utils

import (
	"mime"
	"path/filepath"
	"strconv"
)

func IsImage(path string) bool {
	ext := filepath.Ext(path)
	mimeType := mime.TypeByExtension(ext)
	if len(mimeType) > 4 {
		return mimeType[:5] == "image"
	}
	return len(mimeType) >= 5 && mimeType[:6] == "image"
}

func StringInSlice(s string, list []string) bool {
	for _, str := range list {
		if str == s {
			return true
		}
	}
	return false
}

func StringToFloat32(s string) (float32, error) {
	distance64, err := strconv.ParseFloat(s, 32)
	if err != nil {
		return 0, err
	}
	distance := float32(distance64)
	return distance, nil
}

func AverageVectors(v1, v2 []float32, w1, w2 float32) []float32 {
	result := make([]float32, len(v1))
	for i := range v1 {
		result[i] = (v1[i] * w1) + (v2[i] * w2)
	}
	return result
}
