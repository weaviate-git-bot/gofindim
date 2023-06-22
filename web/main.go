package web

import (
	"context"
	"embed"
	"fmt"
	"net/http"
	"time"

	"github.com/agentx3/gofindim/web/routes"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/weaviate/weaviate-go-client/v4/weaviate"
)

//go:embed react/*
var reactFiles embed.FS

func StartChi(weaviateClient *weaviate.Client) {
	r := chi.NewRouter()

	// A good base middleware stack
	// r.Use(middleware.RequestID)
	// r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(CreateWeaviateMiddleWare(weaviateClient))

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	fileServer := http.FileServer(http.FS(reactFiles))
	// r.NotFound(func(w http.ResponseWriter, r *http.Request) {
	// 	fileServer.ServeHTTP(w, r)
	// })
	apiRouter := chi.NewRouter()
	apiRouter.Use(CreateWeaviateMiddleWare(weaviateClient))
	apiRouter.Use(middleware.Logger)
	apiRouter.HandleFunc("/similar", routes.SimilarHandler)
	apiRouter.HandleFunc("/files", routes.BrowseApiHandler)
	apiRouter.HandleFunc("/scan", routes.ScanHandler)

	r.Mount("/api", apiRouter)
	// r.NotFound(func(w http.ResponseWriter, r *http.Request) {
	// 	http.ServeFile(w, r, "static/react")
	// 	// http.FileServer(http.Dir("static/react")).ServeHTTP(w, r)
	// })
	r.Handle("/static/react/*", http.StripPrefix("/static/react", fileServer))
	// r.Get("/", func(w http.ResponseWriter, r *http.Request) {
	// 	rctx := chi.RouteContext(r.Context())
	// 	serverRoutePrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
	// 	fs := http.StripPrefix(serverRoutePrefix, http.FileServer())
	// 	fs.ServeHTTP(w, r)
	// })
	FileServer(r)
	// fmt.Println(r.URL.Path, r.RequestURI)
	// r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
	// path := "static/react" + r.RequestURI
	// ext := filepath.Ext(r.RequestURI)

	// if _, err := os.Stat(path); os.IsNotExist(err) || ext == "" {
	// 	http.ServeFile(w, r, "static/react/index.html")
	// } else {
	// 	fileServer.ServeHTTP(w, r)
	// }
	// })

	http.ListenAndServe(":8888", r)
}

// FileServer is serving static files.
func FileServer(router *chi.Mux) {
	root := "/"
	fs := http.FileServer(http.Dir(root))
	router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.RequestURI)
		// decodedURI, err := url.PathUnescape(r.RequestURI)
		// if err != nil {
		// 	http.Error(w, "Invalid path", http.StatusBadRequest)
		// 	return
		// }
		// if _, err := os.Stat(root + decodedURI); os.IsNotExist(err) {
		// 	http.StripPrefix(decodedURI, fs).ServeHTTP(w, r)
		// } else if decodedURI == "bundle.js" {
		// http.StripPrefix(decodedURI, fs).ServeHTTP(w, r)
		// } else {
		fs.ServeHTTP(w, r)
		// }
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
