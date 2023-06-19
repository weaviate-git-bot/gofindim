package routes

import (
	"fmt"
	"math"
	"net/http"
	"path"
	"strconv"
	"strings"

	"github.com/agentx3/gofindim/data"
	"github.com/agentx3/gofindim/utils"
	"github.com/gin-gonic/gin"
	"github.com/weaviate/weaviate-go-client/v4/weaviate"
)

type myForm struct {
	DeleteImages []string `form:"delete_images[]"`
}

var fieldsToFetch = []string{"path", "name", "id"}

func SimilarApiHandler(c *gin.Context) {
	weaviateClient := c.MustGet("weaviateClient").(*weaviate.Client)
	var (
		results *[]data.ImageNode
		err     error
	)
	results = &[]data.ImageNode{}
	if weaviateClient == nil {
		c.JSON(500, gin.H{
			"error": "Weaviate client not found in context",
		})
		return
	}
	if c.Request.Method == "POST" {
		results, err = similiarApiPostHandler(c)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

	} else if c.Request.Method == "GET" {
		results, err = similarGetHandler(c)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

	} else {
		c.JSON(404, gin.H{
			"error": "Method not allowed",
		})
	}
	c.JSON(200, gin.H{
		"images": results,
	})
}

func similiarApiPostHandler(c *gin.Context) (*[]data.ImageNode, error) {
	var (
		err     error
		results *[]data.ImageNode
	)
	weaviateClient := c.MustGet("weaviateClient").(*weaviate.Client)
	// Retrieve the text input
	distance, err := utils.StringToFloat32(c.PostForm("distance"))
	if err != nil {
		return nil, err
	}

	text_input := c.PostForm("text_input")

	text_weight, err := utils.StringToFloat32(c.PostForm("text_weight"))
	if err != nil {
		text_weight = 0.5
	}
	image_weight, err := utils.StringToFloat32(c.PostForm("image_weight"))
	if err != nil {
		image_weight = 0.5
	}
	limitStr := c.PostForm("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		return nil, err
	}
	limit = int(math.Max(float64(1), float64(limit)))

	var fakeform myForm
	c.ShouldBind(&fakeform)
	if textInput := c.PostForm("text_input"); textInput != "" {
		if path := c.PostForm("path"); path != "" {
			image, err := data.NewImageFileFromPath(path)
			if err != nil {
				return nil, err
			}
			results, err := data.SearchWeaviateWithTextAndImage(text_input,
				image,
				text_weight,
				image_weight,
				distance,
				limit,
				fieldsToFetch,
				weaviateClient,
			)
			if err != nil {
				return nil, err
			}
			return results, nil
		}
		fmt.Printf("Searching with text %v", textInput)
		if strings.HasPrefix(textInput, "http") {
			imageName := path.Base(textInput)
			imageFile, err := data.NewImageFileFromURL(textInput, imageName)
			if err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
			}
			results, err = data.SearchWeaviateWithImageFile(
				imageFile,
				distance,
				limit,
				fieldsToFetch,
				weaviateClient,
			)
			if err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
			}
		} else {
			fmt.Printf("Searching with text %s", textInput)
			results, err = data.SearchWeaviateWithText(
				textInput,
				distance,
				limit,
				fieldsToFetch,
				weaviateClient,
			)
			if err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
			}
		}

	} else if file, err := c.FormFile("file_input"); err == nil {
		// Retrieve the file from the form data
		results, err = data.SearchWeaviateWithFormFile(file, distance, limit, fieldsToFetch, weaviateClient)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
		}

	} else {
		c.AbortWithError(http.StatusBadRequest, err)
	} // Return a response
	return results, err
}

func similarGetHandler(c *gin.Context) (*[]data.ImageNode, error) {
	path := c.Query("path")
	text_input := c.Query("text_input")

	text_weight, err := utils.StringToFloat32(c.Query("text_weight"))
	if err != nil {
		text_weight = 0.5
	}

	image_weight, err := utils.StringToFloat32(c.Query("image_weight"))
	if err != nil {
		image_weight = 0.5
	}

	distance, err := utils.StringToFloat32(c.Query("distance"))
	if err != nil {
		distance = 0.8
	}

	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		limit = 10
	}

	weaviateClient := c.MustGet("weaviateClient").(*weaviate.Client)
	if path != "" && text_input != "" {
		image, err := data.NewImageFileFromPath(path)
		if err != nil {
			return nil, err
		}
		results, err := data.SearchWeaviateWithTextAndImage(text_input,
			image,
			text_weight,
			image_weight,
			distance,
			limit,
			fieldsToFetch,
			weaviateClient,
		)
		if err != nil {
			return nil, err
		}
		println("Found results for text and image")
		return results, nil

	} else if path != "" {
		results, err := similarPathHandler(path, distance, limit, weaviateClient)
		if err != nil {
			return nil, err
		}
		return results, nil
	}
	return nil, nil
}

func similarPathHandler(
	path string,
	distance float32,
	limit int,
	weaviateClient *weaviate.Client,
) (*[]data.ImageNode, error) {

	results, err := data.SearchWeaviateWithImagePath(path, distance, limit, fieldsToFetch, weaviateClient)
	if err != nil {
		return nil, err
	}
	return results, nil

}

func similarUUIDHandler(c *gin.Context,
	uuid string,
	distance float32,
	limit int,
	weaviateClient *weaviate.Client,
) (*[]data.ImageNode, error) {

	results, err := data.SearchWeaviateWithUUID(uuid, distance, limit, fieldsToFetch, weaviateClient)
	if err != nil {
		return nil, err
	}
	return results, nil

}
