package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ae-gis/suki"
    "github.com/ae-gis/suki/ruuto"
)

const (
	port string = ":9088"
)

func homeHandler(docTitle string) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if len(docTitle) < 1 {
				docTitle = "Default TODO APPS"
			}
			suki.WriteJSON(w, r, map[string]interface{}{
				"code": 200,
				"pesan": docTitle,
			})
	})
}

func main() {
	title := os.Getenv("PESAN")
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	r := ruuto.NewChiRouter()
	r.Use(ruuto.Recovery(), ruuto.Logger())
	
	r.GET("/api/golang", homeHandler(title))

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
