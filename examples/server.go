package main

import (
	"os"

	"github.com/apihub/apihub/api"
	"code.cloudfoundry.org/lager"
)

func main() {
	log := lager.NewLogger("my-server")
	log.RegisterSink(lager.NewWriterSink(os.Stdout, lager.DEBUG))
	server := api.New(log, "tcp", ":8000")
	// server := api.New(log, "unix", "/tmp/apihub.sock")
	server.Start()
	select {}
}
