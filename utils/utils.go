package utils

import (
	"mime"
	"path/filepath"
)

func IsImage(path string) bool {
	ext := filepath.Ext(path)
	mimeType := mime.TypeByExtension(ext)
	if len(mimeType) > 4 {
		return mimeType[:5] == "image"
	}
	return len(mimeType) >= 5 && mimeType[:6] == "image"
}
