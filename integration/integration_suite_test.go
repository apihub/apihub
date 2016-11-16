package integration_test

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"code.cloudfoundry.org/consuladapter/consulrunner"

	"github.com/apihub/apihub"
	"github.com/apihub/apihub/client"
	"github.com/apihub/apihub/client/connection"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"

	"testing"
)

var (
	ApihubAPI     string
	ApihubGateway string
	consulRunner  *consulrunner.ClusterRunner
)

var _ = BeforeSuite(func() {
	var err error
	consulRunner = consulrunner.NewClusterRunner(
		9201+GinkgoParallelNode()*consulrunner.PortOffsetLength,
		1,
		"http",
	)

	consulRunner.Start()
	consulRunner.WaitUntilReady()

	ApihubAPI, err = gexec.Build("github.com/apihub/apihub/cmd/api")
	Expect(err).NotTo(HaveOccurred())

	ApihubGateway, err = gexec.Build("github.com/apihub/apihub/cmd/gateway")
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	gexec.CleanupBuildArtifacts()
	consulRunner.Stop()
})

var _ = BeforeEach(func() {
	Expect(consulRunner.Reset()).To(Succeed())
})

func TestIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
	SetDefaultEventuallyTimeout(time.Second * 10)
	RunSpecs(t, "Integration Suite")
}

func startApihub(network string, addressAPI string, portGateway int) *RunningApihub {
	os.Remove(addressAPI)
	args := []string{"--network", network, "--address", addressAPI, "--consul-server", consulRunner.URL()}

	// Start Apihub Api
	apiSession := runner(exec.Command(ApihubAPI, args...))
	Eventually(apiSession).Should(gbytes.Say("apihub-api.start.started"))

	// Start Apihub Gateway
	args = []string{"--consul-server", consulRunner.URL(), "--port", fmt.Sprintf(":%d", portGateway)}
	gatewaySession := runner(exec.Command(ApihubGateway, args...))
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
	Network        string
	AddressAPI     string
	AddressGateway int
	apihub.Client
	APISession     *gexec.Session
	GatewaySession *gexec.Session
}

func (r *RunningApihub) Stop() error {
	Expect(os.Remove(r.AddressAPI)).To(Succeed())
	if err := r.APISession.Command.Process.Kill(); err != nil {
		return err
	}

	if err := r.GatewaySession.Command.Process.Kill(); err != nil {
		return err
	}

	return nil
}
