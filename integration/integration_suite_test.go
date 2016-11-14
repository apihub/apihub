package integration_test

import (
	"encoding/json"
	"os"
	"os/exec"

	"code.cloudfoundry.org/consuladapter"
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
	consulClient  consuladapter.Client
)

func TestIntegration(t *testing.T) {
	RegisterFailHandler(Fail)

	SynchronizedBeforeSuite(func() []byte {
		var err error
		bins := make(map[string]string)

		bins["api"], err = gexec.Build("github.com/apihub/apihub/cmd/api")
		Expect(err).NotTo(HaveOccurred())

		bins["gateway"], err = gexec.Build("github.com/apihub/apihub/cmd/gateway")
		Expect(err).NotTo(HaveOccurred())

		data, err := json.Marshal(bins)
		Expect(err).NotTo(HaveOccurred())

		consulRunner = consulrunner.NewClusterRunner(
			9101+GinkgoParallelNode()*consulrunner.PortOffsetLength,
			1,
			"http",
		)

		consulRunner.Start()
		consulRunner.WaitUntilReady()

		return data
	}, func(data []byte) {
		bins := make(map[string]string)
		Expect(json.Unmarshal(data, &bins)).To(Succeed())

		ApihubAPI = bins["api"]
		ApihubGateway = bins["gateway"]
	})

	SynchronizedAfterSuite(func() {
	}, func() {
		gexec.CleanupBuildArtifacts()
		consulRunner.Stop()
	})

	RunSpecs(t, "Integration Suite")
}

func startApihub(network string, addressAPI string, portGateway string) *RunningApihub {
	os.Remove(addressAPI)
	args := []string{"--network", network, "--address", addressAPI, "--consul-server", consulRunner.URL()}

	// Start Apihub Api
	apiSession := runner(exec.Command(ApihubAPI, args...))
	Eventually(apiSession).Should(gbytes.Say("started"))

	// Start Apihub Gateway
	args = []string{"--port", portGateway}
	gatewaySession := runner(exec.Command(ApihubGateway, args...))
	Eventually(gatewaySession).Should(gbytes.Say("starting"))

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
	AddressGateway string
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
