package api_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/apihub/apihub"
	"github.com/apihub/apihub/api"
	"github.com/apihub/apihub/client"
	"github.com/apihub/apihub/client/connection"
	"github.com/pivotal-golang/lager/lagertest"

	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/config"
	. "github.com/onsi/gomega"
)

var _ = Describe("The Apihub Server", func() {
	var (
		tmpDir       string
		log          *lagertest.TestLogger
		apihubServer *api.ApihubServer
		apihubClient apihub.Client
	)

	BeforeEach(func() {
		log = lagertest.NewTestLogger("apihub-test")
	})

	AfterEach(func() {
		if tmpDir != "" {
			os.RemoveAll(tmpDir)
		}
	})

	Context("when passed a socket", func() {
		It("listens on the socket provided", func() {
			var err error
			tmpDir, err = ioutil.TempDir(os.TempDir(), "apihub-server-test")
			socketPath := path.Join(tmpDir, "apihub.sock")
			apihubServer = api.New(log, "unix", socketPath)
			Expect(err).NotTo(HaveOccurred())

			err = apihubServer.Start()
			Expect(err).NotTo(HaveOccurred())
			info, err := os.Stat(socketPath)
			Expect(err).NotTo(HaveOccurred())
			Expect(info).NotTo(BeNil())
		})
	})

	Context("when passed a tcp addr", func() {
		It("listens on the address provided", func() {
			var err error
			port := fmt.Sprintf(":%d", 8000+config.GinkgoConfig.ParallelNode)

			apihubServer = api.New(log, "tcp", port)

			err = apihubServer.Start()
			Expect(err).NotTo(HaveOccurred())
			apihubClient = client.New(connection.New("tcp", port))
			Expect(apihubClient.Ping()).To(Succeed())
		})
	})

})
