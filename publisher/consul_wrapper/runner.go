package consul_wrapper

import (
	"io/ioutil"
	"os/exec"
	"strconv"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

type Runner struct {
	Address string
	Port    int
	Scheme  string
	session *gexec.Session
}

func NewRunner(scheme string, address string, port int) *Runner {
	return &Runner{
		Scheme:  scheme,
		Address: address,
		Port:    port,
	}
}

func (c *Runner) Start() error {
	dataDir, err := ioutil.TempDir("", "consul_data_dir")
	Expect(err).NotTo(HaveOccurred())

	args := []string{"agent", "-server", "-data-dir", dataDir, "-advertise", c.Address, "-http-port", strconv.Itoa(c.Port), "-bootstrap-expect", "1"}

	cmd := exec.Command("consul", args...)
	session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())
	Eventually(session, "15s").Should(gbytes.Say("agent: Synced service 'consul'"))
	return nil
}

func (c *Runner) Stop() error {
	return c.session.Command.Process.Kill()
}
