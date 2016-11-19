package test_helpers

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/apihub/apihub"
	"github.com/apihub/apihub/client"
	"github.com/apihub/apihub/client/connection"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

func BuildArtifacts() (string, string) {
	apiBin, err := gexec.Build("github.com/apihub/apihub/cmd/api")
	Expect(err).NotTo(HaveOccurred())

	gatewayBin, err := gexec.Build("github.com/apihub/apihub/cmd/gateway")
	Expect(err).NotTo(HaveOccurred())
	return apiBin, gatewayBin
}

func CleanArtifacts() {
	gexec.CleanupBuildArtifacts()
}

func StartApihub(apiBin string, gatewayBin string, network string, addressAPI string, portGateway int, consulURL string) *RunningApihub {
	os.Remove(addressAPI)
	args := []string{"--network", network, "--address", addressAPI, "--consul-server", consulURL}

	// Start Apihub Api
	apiSession := runner(exec.Command(apiBin, args...))
	Eventually(apiSession).Should(gbytes.Say("apihub-api.start.started"))

	// Start Apihub Gateway
	args = []string{"--consul-server", consulURL, "--port", fmt.Sprintf(":%d", portGateway)}
	gatewaySession := runner(exec.Command(gatewayBin, args...))
	Eventually(gatewaySession).Should(gbytes.Say("apihub-gateway.start.starting"))

	ah := &RunningApihub{
		Network:        network,
		AddressAPI:     addressAPI,
		AddressGateway: portGateway,
		Client:         client.New(connection.New(network, addressAPI)),
		APISession:     apiSession,
		GatewaySession: gatewaySession,
	}

	return ah
}

func runner(cmd *exec.Cmd) *gexec.Session {
	session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())
	return session
}

type RunningApihub struct {
	apihub.Client

	Network        string
	AddressAPI     string
	AddressGateway int
	APISession     *gexec.Session
	GatewaySession *gexec.Session
}

func (r *RunningApihub) Stop() error {
	os.Remove(r.AddressAPI)
	if err := r.APISession.Command.Process.Kill(); err != nil {
		return err
	}

	if err := r.GatewaySession.Command.Process.Kill(); err != nil {
		return err
	}

	return nil
}
