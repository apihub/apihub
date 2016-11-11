package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/apihub/apihub/apihubfakes"
	"github.com/apihub/apihub/gateway"

	"code.cloudfoundry.org/lager"
)

var (
	port = flag.String("port", ":8080", "Port to be used")
)

func main() {
	flag.Parse()

	// Configure log
	logger := lager.NewLogger("apihub-gateway")
	logger.RegisterSink(lager.NewWriterSink(os.Stdout, lager.DEBUG))

	subscriber := new(apihubfakes.FakeServiceSubscriber)
	// Configure and start server
	reverseProxyCreator := gateway.NewReverseProxyCreator()
	gw := gateway.New(*port, subscriber, reverseProxyCreator)

	if err := gw.Start(logger); err != nil {
		panic(fmt.Errorf("Failed to start Apihub Gateway: `%s`.", err))
	}
}
