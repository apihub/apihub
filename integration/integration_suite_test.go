package integration_test

import (
	"encoding/json"
	"io"
	"os"
	"os/exec"

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
	ApihubAPI string
)

func TestIntegration(t *testing.T) {
	RegisterFailHandler(Fail)

	SynchronizedBeforeSuite(func() []byte {
		var err error
		bins := make(map[string]string)

		bins["api"], err = gexec.Build("github.com/apihub/apihub/cmd/api")
		Expect(err).NotTo(HaveOccurred())

		data, err := json.Marshal(bins)
		Expect(err).NotTo(HaveOccurred())

		return data
	}, func(data []byte) {
		bins := make(map[string]string)
		Expect(json.Unmarshal(data, &bins)).To(Succeed())

		ApihubAPI = bins["api"]
	})

	RunSpecs(t, "Integration Suite")
}

func startApihub(network string, address string) *RunningApihub {
	args := []string{"--network", network, "--address", address}

	buf := gbytes.NewBuffer()
	cmd := exec.Command(ApihubAPI, args...)
	session, err := gexec.Start(cmd, io.MultiWriter(buf, GinkgoWriter), GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())
	Eventually(session).Should(gbytes.Say("apihub-started"))
	go func() { session.Wait() }()

	ah := &RunningApihub{
		Network: network,
		Address: address,
		Client:  client.New(connection.New(network, address)),
		Session: session,
	}

	return ah
}

type RunningApihub struct {
	Network string
	Address string
	apihub.Client
	Session *gexec.Session
}

func (r *RunningApihub) Stop() error {
	Expect(os.Remove(r.Address)).To(Succeed())
	return r.Session.Command.Process.Kill()
}
