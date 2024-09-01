package main

import "github.com/ZygmuntJakub/mkino-extension/internal/server"

func main() {
	srv := server.NewServer()
	srv.Run()
}
