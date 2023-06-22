package routes

import (
	"fmt"
	"net/http"
	"os"

	"github.com/agentx3/gofindim/data"
	"github.com/weaviate/weaviate-go-client/v4/weaviate"
)

func ScanHandler(w http.ResponseWriter, r *http.Request) {
	weaviateClient := r.Context().Value("weaviateClient").(*weaviate.Client)

	if r.Method == "POST" {
		path := r.URL.Query().Get("path")
		if path == "" {
			path = r.FormValue("path")
			if path == "" {
				http.Error(w, "Path not found", http.StatusBadRequest)
				return
			}
		}
		go func() {
			fmt.Printf("Scanning\n")
			if fileInfo, err := os.Stat(path); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			} else if !fileInfo.IsDir() {
				img, err := data.NewImageFileFromPath(path)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				err = data.InsertIntoWeaviate(img, weaviateClient)
			} else if fileInfo.IsDir() {
				err := data.InsertDirectoryIntoWeaviate(path, weaviateClient)
				if err != nil {
					fmt.Println(err)
					return
				}
			}
		}()
		http.StatusText(http.StatusOK)
	}

}
