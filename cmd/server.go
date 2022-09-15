package cmd

import "github.com/zx06/saut/internal/server"

func RunServer() {
	httpServer := server.NewHTTPServer()
	httpServer.Start()
}
