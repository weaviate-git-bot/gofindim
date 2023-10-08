package web

import (
	"context"
	"embed"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/agentx3/gofindim/web/routes"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/weaviate/weaviate-go-client/v4/weaviate"
)

//go:embed frontend/dist/*
var reactFiles embed.FS

func customFileServer(fs http.FileSystem) http.Handler {
	fileServer := http.FileServer(fs)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		strippedPath := strings.TrimPrefix(r.URL.Path, "/web")
		fmt.Println("Path", r.URL.Path)
		r.URL.Path = "/frontend/dist" + strippedPath
		fmt.Println("Path", r.URL.Path)
		fileServer.ServeHTTP(w, r)
	})
}

func StartChi(weaviateClient *weaviate.Client) {
	r := chi.NewRouter()

	// A good base middleware stack
	// r.Use(middleware.RequestID)
	// r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(CreateWeaviateMiddleWare(weaviateClient))
	r.Use(CorsMiddleware())

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	apiRouter := chi.NewRouter()
	apiRouter.Use(CreateWeaviateMiddleWare(weaviateClient))
	apiRouter.Use(CorsMiddleware())
	apiRouter.Use(middleware.Logger)
	apiRouter.HandleFunc("/similar", routes.SimilarHandler)
	apiRouter.HandleFunc("/files", routes.BrowseApiHandler)
	apiRouter.HandleFunc("/scan", routes.ScanHandler)

	r.Mount("/api", apiRouter)
	r.Handle("/web/*", customFileServer(http.FS(reactFiles)))
	FileServer(r)

	http.ListenAndServe(":8888", r)
}

// FileServer is serving static files.
func FileServer(router *chi.Mux) {
	root := "/"
	fs := http.FileServer(http.Dir(root))
	router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.RequestURI, "router")
		fs.ServeHTTP(w, r)
	})
}
func CreateWeaviateMiddleWare(weaviateClient *weaviate.Client) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), "weaviateClient", weaviateClient)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func CorsMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Access-Control-Allow-Origin", "http://localhost:*")
			next.ServeHTTP(w, r)
		})
	}
}
