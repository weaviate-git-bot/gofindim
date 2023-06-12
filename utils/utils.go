package utils

import (
	"fmt"
	"mime"
)

func IsImage(filepath string) bool {
	mimeType := mime.TypeByExtension(filepath)
	if len(mimeType) > 4 {
		fmt.Println(mimeType[:5])
		return mimeType[:5] == "image"
	}
	return len(mimeType) >= 5 && mimeType[:6] == "image"
}
