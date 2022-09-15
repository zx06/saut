package server

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zx06/saut/internal/handler"
)

type HTTPServer struct {
	server *http.Server
}

func registerRouters(srv *HTTPServer) {
	engine := gin.New()
	srv.server.Handler = engine
	engine.Use(
		gin.Logger(),
		gin.Recovery(),
	)
	wsGroup := engine.Group("/ws")
	{
		wsGroup.GET("/terminal", handler.TerminalHandler)
	}
}

func NewHTTPServer() *HTTPServer {
	return &HTTPServer{
		server: &http.Server{
			Addr: ":12345",
		},
	}
}

func (srv *HTTPServer) Start() {
	quit := make(chan os.Signal, 1)
	registerRouters(srv)
	go func() {
		if err := srv.server.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Printf("listen: %s\n", err)
		}
	}()
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
	log.Println("Server exiting")
}

func (srv *HTTPServer) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
	log.Println("Server exiting")
}
