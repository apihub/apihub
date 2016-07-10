package main

import (
	"os"

	"github.com/apihub/apihub/api"
	"github.com/pivotal-golang/lager"
)

func main() {
	log := lager.NewLogger("my-server")
	log.RegisterSink(lager.NewWriterSink(os.Stdout, lager.DEBUG))
	server := api.New(log, "tcp", ":8000")
	server.Start()
	select {}
}
