package api_test

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/apihub/apihub"
	"github.com/apihub/apihub/api"
	"github.com/apihub/apihub/client"
	"github.com/apihub/apihub/client/connection"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"code.cloudfoundry.org/lager/lagertest"
)

var _ = Describe("When a client connects", func() {
	var (
		apihubClient apihub.Client
		apihubServer *api.ApihubServer
		err          error
		log          *lagertest.TestLogger
		tmpDir       string
	)

	BeforeEach(func() {
		log = lagertest.NewTestLogger("apihub-test")
		tmpDir, err = ioutil.TempDir(os.TempDir(), "apihub-server-test")
		socketPath := path.Join(tmpDir, "apihub.sock")
		apihubServer = api.New(log, "unix", socketPath)
		Expect(err).NotTo(HaveOccurred())

		Expect(apihubServer.Start()).NotTo(HaveOccurred())
		apihubClient = client.New(connection.New("unix", socketPath))
	})

	AfterEach(func() {
		if tmpDir != "" {
			os.RemoveAll(tmpDir)
		}
	})

	Describe("and sends a Ping Request", func() {
		Context("and the server is up and running", func() {
			It("does not return an error", func() {
				Expect(apihubClient.Ping()).NotTo(HaveOccurred())
			})
		})

		Context("when the server is not up and running", func() {
			BeforeEach(func() {
				Expect(apihubServer.Stop()).NotTo(HaveOccurred())
			})

			It("returns an error", func() {
				Expect(apihubClient.Ping()).To(HaveOccurred())
			})
		})
	})

})
