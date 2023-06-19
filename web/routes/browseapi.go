package routes

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nfnt/resize"
)

func BrowseApiHandler(c *gin.Context) {

	if c.Request.Method == "GET" {
		hidden := false
		path := c.Query("path")
		if path == "" {
			path = "/"
		}
		hiddenStr := c.Query("hidden")
		if hiddenStr == "true" {
			hidden = true
		}
		files, err := ioutil.ReadDir(path)
		if err != nil {
			c.AbortWithStatus(500)
		}
		var fileInfoList []gin.H
		for _, file := range files {
			if strings.HasPrefix(file.Name(), ".") && !hidden {
				continue
			}
			var thumbnailPath string
			if isMediaFile(file.Name()) {
				thumbnailPath, err = generateAndFetchThumbnail(path, path, file.Name())
				if err != nil {
					c.AbortWithStatus(500)
					return
				}
			}
			fileInfoList = append(fileInfoList, gin.H{
				"name":         file.Name(),
				"isDir":        file.IsDir(),
				"thumbnailUrl": thumbnailPath,
				"path":         filepath.Join(path, file.Name()),
			})
		}
		c.JSON(http.StatusOK, gin.H{"files": fileInfoList})

	}

}

func isMediaFile(file string) bool {
	mimeType := mime.TypeByExtension(filepath.Ext(file))
	if strings.Contains(mimeType, "svg") {
		return false
	}
	return strings.HasPrefix(mimeType, "image/")
}

func generateAndFetchThumbnail(path, dir, file string) (string, error) {
	imagePath := filepath.Join(dir, file)
	thumbnailDir := filepath.Join("/", "mnt", "ramdisk", "thumbnails")
	thumbnailPath := filepath.Join(thumbnailDir, file)

	if _, err := os.Stat(thumbnailPath); os.IsNotExist(err) {
		if err := os.MkdirAll(thumbnailDir, 0755); err != nil {
			return "", err
		}
		imageFile, err := os.Open(filepath.Join(dir, file))
		if err != nil {
			return "", err
		}
		defer imageFile.Close()
		thumbnailFile, err := os.Create(thumbnailPath)
		if err != nil {
			return "", err
		}
		defer thumbnailFile.Close()
		decodedImage, _, err := image.Decode(imageFile)
		if err != nil {
			return "", fmt.Errorf("error decoding image: %v", err)
		}
		thumb := resize.Thumbnail(100, 100, decodedImage, resize.Lanczos3)
		switch filepath.Ext(imagePath) {
		case ".jpg", ".jpeg":
			return thumbnailPath, jpeg.Encode(thumbnailFile, thumb, nil)
		case ".png":
			return thumbnailPath, png.Encode(thumbnailFile, thumb)
		default:
			return thumbnailPath, jpeg.Encode(thumbnailFile, thumb, nil)
		}
	}
	return thumbnailPath, nil
}
