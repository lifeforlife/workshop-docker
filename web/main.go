package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"strings"
	"flag"

	"github.com/ae-gis/suki"
    "github.com/ae-gis/suki/ruuto"
	"github.com/thedevsaddam/renderer"
)

var rendr *renderer.Render

const (
	port string = ":9099"
	STATIC_DIR string = "/static"
)

func init() {
	rendr = renderer.New()
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if err := rendr.Template(w, http.StatusOK, []string{"/static/home.tpl"}, nil); err != nil {
		suki.Error("Home Handler",
			suki.Field("Error", err),
		)
	}
}

// FileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
func FileServer(r ruuto.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit URL parameters.")
	}

	fs := http.StripPrefix(path, http.FileServer(root))

	if path != "/" && path[len(path)-1] != '/' {
		r.GET(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.GET(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))
}

func main() {
	directory := flag.String("d", STATIC_DIR, "the directory of static file to host")
	flag.Parse()

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	r := ruuto.NewChiRouter()
	r.Use(ruuto.Recovery(), ruuto.Logger())

	

	FileServer(r, STATIC_DIR, http.Dir(*directory))
	
	r.GET("/", homeHandler)

	srv := &http.Server{
			Addr: port,
			Handler: r,
			ReadTimeout:  60 * time.Second,
			WriteTimeout: 60 * time.Second,
			IdleTimeout:  60 * time.Second,
	}
	go func() {
		suki.Info("Listening on port ",
			suki.Field("Port", port),
		)
		if err := srv.ListenAndServe(); err != nil {
			suki.Error("listen",
				suki.Field("Error", err),
			)
		}
	}()

	<-stopChan
	suki.Info("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	if err := srv.Shutdown(ctx); err != nil {
		suki.Error("Shutdown => ",
			suki.Field("Error", err),
		)
	}
	defer cancel()
	suki.Info("Server gracefully stopped!")
}


