package main

import (
	"flag"
	"os"

	"code.cloudfoundry.org/lager"
	"github.com/apihub/apihub/api"
	"github.com/apihub/apihub/storage"
)

var (
	network = flag.String("network", "unix", "Either `tcp` or `unix`")
	address = flag.String("address", "/tmp/apihub.sock", "Port for `tcp` or filepath for `unix`")
)

func main() {
	flag.Parse()

	// Configure log
	log := lager.NewLogger("apihub")
	log.RegisterSink(lager.NewWriterSink(os.Stdout, lager.DEBUG))

	// Start api server
	store := storage.New()
	server := api.New(log, *network, *address, store)
	server.Start()
	//TODO: Remove this!
	select {}
}
