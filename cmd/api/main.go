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
	logger := lager.NewLogger("apihub-api")
	logger.RegisterSink(lager.NewWriterSink(os.Stdout, lager.DEBUG))

	// Configure and start server
	store := storage.New()
	server := api.New(logger, *network, *address, store)
	server.Start(true)
}
