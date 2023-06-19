package web

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/agentx3/gofindim/web/routes"
	"github.com/gin-gonic/gin"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/weaviate/weaviate-go-client/v4/weaviate"
)

func Start(weaviateClient *weaviate.Client) {
	r := gin.Default()
	exePath, err := os.Executable()
	if err != nil {
		panic(err)
	}
	dir := filepath.Dir(exePath)
	r.LoadHTMLGlob(filepath.Join(dir, "web", "templates", "*.html"))
	r.LoadHTMLGlob(filepath.Join(dir, "web", "templates", "components", "*.html"))
	r.Use(func(c *gin.Context) {
		// Inject the Weaviate client into the context
		c.Set("weaviateClient", weaviateClient)
		// Continue processing the request
		c.Next()
	})
	r.Static("/static", filepath.Join(dir, "static"))
	r.Static("/thumbnails", filepath.Join("mnt", "ramdisk", "thumbnails"))
	// r.Static("/static", filepath.Join(dir, "web", "static", "css"))
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	api := r.Group("/api")
	api.Any("similar", routes.SimilarApiHandler)
	api.Any("files", routes.BrowseApiHandler)
	r.NoRoute(func(c *gin.Context) {
		c.File(filepath.Join(dir, "static", "react", "index.html"))
	})
	// r.GET("/", func(c *gin.Context) {
	// 	routes.RootHandler(c)
	// })

	r.Run(":8888") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func CreateWeaviateMiddleWare(weaviateClient *weaviate.Client) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), "weaviateClient", weaviateClient)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func WeaviateMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// create new context from `r` request context, and assign key `"user"`
		// to value of `"123"`
		ctx := context.WithValue(r.Context(), "user", "123")

		// call the next handler in the chain, passing the response writer and
		// the updated request object with the new context value.
		//
		// note: context.Context values are nested, so any previously set
		// values will be accessible as well, and the new `"user"` key
		// will be accessible from this point forward.
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
func StartChi(weaviateClient *weaviate.Client) {
	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(CreateWeaviateMiddleWare(weaviateClient))

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hi"))
	})

	// RESTy routes for "articles" resource
	// r.Route("/articles", func(r chi.Router) {
	// 	r.With(paginate).Get("/", listArticles)                           // GET /articles
	// 	r.With(paginate).Get("/{month}-{day}-{year}", listArticlesByDate) // GET /articles/01-16-2017

	// 	r.Post("/", createArticle)       // POST /articles
	// 	r.Get("/search", searchArticles) // GET /articles/search

	// 	// Regexp url parameters:
	// 	r.Get("/{articleSlug:[a-z-]+}", getArticleBySlug) // GET /articles/home-is-toronto

	// 	// Subrouters:
	// 	r.Route("/{articleID}", func(r chi.Router) {
	// 		r.Use(ArticleCtx)
	// 		r.Get("/", getArticle)       // GET /articles/123
	// 		r.Put("/", updateArticle)    // PUT /articles/123
	// 		r.Delete("/", deleteArticle) // DELETE /articles/123
	// 	})
	// })

	// // Mount the admin sub-router
	// r.Mount("/admin", adminRouter())

	http.ListenAndServe(":9999", r)
}
