package integration_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"code.cloudfoundry.org/lager/lagertest"
	"github.com/apihub/apihub/gateway"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Gateway", func() {
	var (
		port                string
		gw                  *gateway.Gateway
		logger              *lagertest.TestLogger
		reverseProxyCreator gateway.ReverseProxyCreator
		spec                gateway.ReverseProxySpec
	)

	BeforeEach(func() {
		logger = lagertest.NewTestLogger("gateway")
		port = fmt.Sprintf(":909%d", GinkgoParallelNode())
		reverseProxyCreator = gateway.NewReverseProxyCreator()

		gw = gateway.New(port, reverseProxyCreator)

		spec = gateway.ReverseProxySpec{
			Handle:   "my-handle",
			Backends: []string{},
		}
	})

	Describe("AddService", func() {
		It("proxies the request to the service backend", func() {
			backendServer := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
				rw.Write([]byte("Hello World."))
			}))
			defer backendServer.Close()

			spec.Backends = []string{fmt.Sprintf("http://%s", backendServer.Listener.Addr().String())}

			Expect(gw.AddService(logger, spec)).To(Succeed())

			req, err := http.NewRequest(http.MethodGet, "http://my-handle.apihub.dev", nil)
			Expect(err).NotTo(HaveOccurred())

			var rw *httptest.ResponseRecorder
			rw = httptest.NewRecorder()
			gw.ServeHTTP(rw, req)

			Expect(rw.Body.String()).To(Equal("Hello World."))
		})
	})
})
