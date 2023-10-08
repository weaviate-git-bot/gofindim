package routes

import (
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"math"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/nfnt/resize"
)

func BrowseApiHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {
		hidden := false
		path := r.URL.Query().Get("path")
		if path == "" {
			path = "/"
		}

		hiddenStr := r.URL.Query().Get("hidden")
		if hiddenStr == "true" {
			hidden = true
		}
		limitStr := r.URL.Query().Get("limit")
		if limitStr == "" {
			limitStr = "20"
		}
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		pageStr := r.URL.Query().Get("page")
		if pageStr == "" {
			pageStr = "0"
		}
		page, err := strconv.Atoi(pageStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Println("Trying to access", path)
		files, err := ioutil.ReadDir(path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		var fileInfoList []map[string]interface{}
		for _, file := range files {
			if strings.HasPrefix(file.Name(), ".") && !hidden {
				continue
			}
			var thumbnailPath string
			if isMediaFile(file.Name()) {
				thumbnailPath, err = generateAndFetchThumbnail(path, path, file.Name())
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}
			fileInfoList = append(fileInfoList, map[string]interface{}{
				"name":         file.Name(),
				"isDir":        file.IsDir(),
				"thumbnailUrl": thumbnailPath,
				"path":         filepath.Join(path, file.Name()),
			})
		}
		r.Header.Set("Content-Type", "application/json")
		itemCount := len(fileInfoList)
		lastIndex := math.Min(float64(page*limit+limit), float64(len(fileInfoList)))
		json.NewEncoder(w).Encode(map[string]interface{}{"files": fileInfoList[page*limit : int(lastIndex)], "itemCount": itemCount})

	}

}

func isMediaFile(file string) bool {
	mimeType := mime.TypeByExtension(filepath.Ext(file))
	if strings.Contains(mimeType, "svg") {
		return false
	}
	return strings.HasPrefix(mimeType, "image/") || strings.HasPrefix(mimeType, "video/")
}
func isVideoFile(file string) bool {
	mimeType := mime.TypeByExtension(filepath.Ext(file))
	return strings.HasPrefix(mimeType, "video/")
}

func generateAndFetchThumbnail(path, dir, file string) (string, error) {
	imagePath := filepath.Join(dir, file)
	// thumbnailDir := filepath.Join("/", "mnt", "ramdisk", "thumbnails", "large")
	thumbnailDir := filepath.Join(path, ".thumbnails")
	filePng := file + ".png"
	inDirThumbPath := filepath.Join(path, ".thumbnails", filePng)
	// thumbnailPath := filepath.Join(thumbnailDir, filePng)
	thumbnailPath := inDirThumbPath

	if _, err := os.Stat(inDirThumbPath); !os.IsNotExist(err) {
		return inDirThumbPath, nil
	} else if isVideoFile(file) {
		return "", nil
	}
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
			return "", fmt.Errorf("error decoding image %v: %v", imagePath, err)
		}
		thumb := resize.Thumbnail(100, 100, decodedImage, resize.Lanczos3)
		switch filepath.Ext(imagePath) {
		// case ".jpg", ".jpeg":
		// 	return thumbnailPath, jpeg.Encode(thumbnailFile, thumb, nil)
		case ".png":
			return thumbnailPath, png.Encode(thumbnailFile, thumb)
		default:
			return thumbnailPath, png.Encode(thumbnailFile, thumb)
		}
	}
	return thumbnailPath, nil
}
