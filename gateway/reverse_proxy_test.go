package gateway_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"code.cloudfoundry.org/lager/lagertest"

	"github.com/apihub/apihub/gateway"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ReverseProxyCreator", func() {

	var (
		backendServer *httptest.Server
		creator       gateway.ReverseProxyCreator
		spec          gateway.ReverseProxySpec
		logger        *lagertest.TestLogger
	)

	JustBeforeEach(func() {
		backendServer = httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			Expect(req.URL.Query()["foo"]).To(Equal([]string{"bar"}))
			Expect(req.URL.Query()["bar"]).To(Equal([]string{"foo"}))
			rw.Write([]byte("Hello world."))
		}))

		spec = gateway.ReverseProxySpec{
			Host:     "my-host",
			Backends: []string{fmt.Sprintf("http://%s", backendServer.Listener.Addr().String())},
		}

		logger = lagertest.NewTestLogger("hostr-provider")
		creator = gateway.NewReverseProxyCreator()
	})

	AfterEach(func() {
		backendServer.Close()
	})

	Describe("Create", func() {
		It("creates a reverse proxy", func() {
			reverseProxy, err := creator.Create(logger, spec)
			Expect(err).NotTo(HaveOccurred())

			req, err := http.NewRequest(http.MethodGet, "http://my-host.apihub.dev?foo=bar&bar=foo", nil)
			Expect(err).NotTo(HaveOccurred())

			rw := httptest.NewRecorder()
			reverseProxy.ServeHTTP(rw, req)

			Expect(rw.Body.String()).To(Equal("Hello world."))
			Expect(rw.Header().Get("Via")).NotTo(BeEmpty())
		})

		Context("when fails to create a proxy", func() {
			var badSpec gateway.ReverseProxySpec

			BeforeEach(func() {
				badSpec = gateway.ReverseProxySpec{
					Host:     "my-host",
					Backends: []string{},
				}
			})

			It("returns an error", func() {
				_, err := creator.Create(logger, badSpec)
				Expect(err).To(MatchError(ContainSubstring("Backends cannot be empty")))
			})
		})
	})

})
