package main

import (
	"github.com/ZygmuntJakub/mkino-extension/internal/collector"
	"github.com/ZygmuntJakub/mkino-extension/internal/server"
)

func main() {
	go func() {
		collector.RunCollector()
	}()
	srv := server.NewServer()
	srv.Run()
}
