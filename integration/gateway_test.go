package integration_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	"code.cloudfoundry.org/lager/lagertest"
	"github.com/apihub/apihub/gateway"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Gateway", func() {
	var (
		portGateway         int
		gw                  *gateway.Gateway
		logger              *lagertest.TestLogger
		reverseProxyCreator gateway.ReverseProxyCreator
		spec                gateway.ReverseProxySpec
	)

	BeforeEach(func() {
		logger = lagertest.NewTestLogger("gateway")
		portGateway = 9000 + GinkgoParallelNode()
		reverseProxyCreator = gateway.NewReverseProxyCreator()

		gw = gateway.New(fmt.Sprintf(":%d", portGateway), reverseProxyCreator)

		spec = gateway.ReverseProxySpec{
			Host:     "my-host.apihub.dev",
			Backends: []string{"http://server-a"},
		}
	})

	Describe("Adding a service", func() {
		It("adds a service", func() {
			Expect(gw.AddService(logger, spec)).To(Succeed())
		})
	})

	Describe("Proxing requests through the Gateway", func() {
		It("proxies the request to the service backend", func() {
			backendServer := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
				rw.Write([]byte("Hello World."))
			}))
			defer backendServer.Close()

			spec.Backends = []string{fmt.Sprintf("http://%s", backendServer.Listener.Addr().String())}

			Expect(gw.AddService(logger, spec)).To(Succeed())

			req, err := http.NewRequest(http.MethodGet, "http://my-host.apihub.dev", nil)
			Expect(err).NotTo(HaveOccurred())

			rw := httptest.NewRecorder()
			gw.ServeHTTP(rw, req)

			Expect(rw.Body.String()).To(Equal("Hello World."))
		})

		Context("when the service does not respond as expected", func() {
			BeforeEach(func() {
				Expect(gw.AddService(logger, spec)).To(Succeed())
			})

			It("returns a bad gateway error", func() {
				req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://127.0.0.1:%d", portGateway), nil)
				Expect(err).NotTo(HaveOccurred())
				req.Host = spec.Host

				rw := httptest.NewRecorder()
				gw.ServeHTTP(rw, req)

				Expect(rw.Code).To(Equal(http.StatusBadGateway))
			})
		})

		Context("when the service timeout is reached", func() {
			var backendServer *httptest.Server

			BeforeEach(func() {
				backendServer = httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
					select {
					case <-time.After(11 * time.Millisecond):
					}
				}))
				spec.Timeout = time.Millisecond * 10
				spec.Backends = []string{fmt.Sprintf("http://%s", backendServer.Listener.Addr().String())}
				Expect(gw.AddService(logger, spec)).To(Succeed())
			})

			AfterEach(func() {
				backendServer.Close()
			})

			It("interrupts the request and returns gateway timeout", func() {
				req, err := http.NewRequest(http.MethodGet, "http://my-host.apihub.dev", nil)
				Expect(err).NotTo(HaveOccurred())

				rw := httptest.NewRecorder()
				gw.ServeHTTP(rw, req)

				Expect(rw.Code).To(Equal(http.StatusGatewayTimeout))
				Expect(rw.Body.String()).To(ContainSubstring("i/o timeout"))
			})
		})
	})

	Describe("Removing a service", func() {
		var backendServer *httptest.Server

		BeforeEach(func() {
			backendServer = httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
				rw.Write([]byte("Hello World."))
			}))
			spec.Backends = []string{fmt.Sprintf("http://%s", backendServer.Listener.Addr().String())}
			Expect(gw.AddService(logger, spec)).To(Succeed())
		})

		AfterEach(func() {
			backendServer.Close()
		})

		It("removes a service", func() {
			Expect(gw.RemoveService(logger, spec.Host)).To(Succeed())
		})

		It("stops proxying requests to the service", func() {
			req, err := http.NewRequest(http.MethodGet, "http://my-host.apihub.dev", nil)
			Expect(err).NotTo(HaveOccurred())

			var rw *httptest.ResponseRecorder
			rw = httptest.NewRecorder()
			gw.ServeHTTP(rw, req)
			Expect(rw.Code).To(Equal(http.StatusOK))

			Expect(gw.RemoveService(logger, spec.Host)).To(Succeed())

			Eventually(func() int {
				rw = httptest.NewRecorder()
				gw.ServeHTTP(rw, req)
				return rw.Code
			}).Should(Equal(http.StatusNotFound))
		})

		Context("when service is not found", func() {
			It("returns an error", func() {
				Expect(gw.RemoveService(logger, "invalid-host")).To(MatchError(ContainSubstring("service not found: 'invalid-host'")))
			})
		})
	})
})
