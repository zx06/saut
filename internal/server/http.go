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

	"github.com/zx06/saut/internal/pkg/auth"

	"github.com/gin-gonic/gin"
	"github.com/zx06/saut/internal/handler"
)

type HTTPServer struct {
	server *http.Server
}

func registerRouters(srv *HTTPServer, auth *auth.Auth) {
	engine := gin.New()
	srv.server.Handler = engine
	engine.Use(
		gin.Logger(),
		gin.Recovery(),
	)
	authGroup := engine.Group("/auth")
	{
		authGroup.POST("/login", auth.LoginHandler())
	}
	assetsGroup := engine.Group("/assets", auth.AuthMiddleware())
	{
		assetsGroup.GET("/", func(c *gin.Context) {

		})
		assetsGroup.GET("/:id", func(c *gin.Context) {

		})
		assetsGroup.POST("/", func(c *gin.Context) {

		})
		assetsGroup.PUT("/:id", func(c *gin.Context) {

		})
	}
	wsGroup := engine.Group("/ws", auth.AuthMiddleware())
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
	auth := auth.NewAuth()
	registerRouters(srv, auth)
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
