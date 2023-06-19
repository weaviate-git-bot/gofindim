package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/weaviate/weaviate-go-client/v4/weaviate"
)

func RootHandler(c *gin.Context) {
	weaviateClient := c.MustGet("weaviateClient").(*weaviate.Client)
	if weaviateClient == nil {
		c.JSON(500, gin.H{
			"error": "Weaviate client not found in context",
		})
		return
	}
	c.HTML(200, "base.html", gin.H{
		"Title":      "My Website",
		"Header":     "Welcome to My Website",
		"Subcontent": "This is a sub-content.",
	})

	// c.HTML(200, "root.html", gin.H{
	// 	"Title": "<div style={background-color: blue;}>GoFindim</div>",
	// 	"Names": []string{"Alice", "Bob", "<script>alert('you have been pwned')</script>"},
	// })
}
