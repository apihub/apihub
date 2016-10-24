package api_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"code.cloudfoundry.org/lager/lagertest"
	"github.com/apihub/apihub"
	"github.com/apihub/apihub/api"
	"github.com/apihub/apihub/apihubfakes"
	"github.com/apihub/apihub/client"
	"github.com/apihub/apihub/client/connection"

	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/config"
	. "github.com/onsi/gomega"
)

var _ = Describe("Apihub Server", func() {
	var (
		tmpDir       string
		log          *lagertest.TestLogger
		storage      apihub.Storage
		apihubServer *api.ApihubServer
		apihubClient apihub.Client
	)

	BeforeEach(func() {
		log = lagertest.NewTestLogger("apihub-test")
		storage = new(apihubfakes.FakeStorage)
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
			apihubServer = api.New(log, "unix", socketPath, storage)
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

			apihubServer = api.New(log, "tcp", port, storage)

			err = apihubServer.Start()
			Expect(err).NotTo(HaveOccurred())
			apihubClient = client.New(connection.New("tcp", port))
			Expect(apihubClient.Ping()).To(Succeed())
		})
	})
})
