package gateway_test

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"

	"code.cloudfoundry.org/lager/lagertest"
	"github.com/apihub/apihub/gateway"
	"github.com/apihub/apihub/gateway/gatewayfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Gateway", func() {
	var (
		port                    string
		gw                      *gateway.Gateway
		logger                  *lagertest.TestLogger
		fakeReverseProxyCreator *gatewayfakes.FakeReverseProxyCreator
		fakeReverseProxy        *gatewayfakes.FakeReverseProxy
	)

	BeforeEach(func() {
		logger = lagertest.NewTestLogger("test-gateway")
		port = fmt.Sprintf(":909%d", GinkgoParallelNode())
		fakeReverseProxyCreator = new(gatewayfakes.FakeReverseProxyCreator)
		fakeReverseProxy = new(gatewayfakes.FakeReverseProxy)

		gw = gateway.New(port, fakeReverseProxyCreator)
	})

	var spec gateway.ReverseProxySpec

	BeforeEach(func() {
		spec = gateway.ReverseProxySpec{
			Host:     "my-host.apihub.dev",
			Backends: []string{},
		}

		fakeReverseProxyCreator.CreateReturns(fakeReverseProxy, nil)
	})

	Describe("Stop", func() {
		var (
			client *http.Client
		)

		BeforeEach(func() {
			go gw.Start(logger)
			client = &http.Client{}
		})

		It("stops accepting new connections", func() {
			Eventually(func() error {
				_, err := http.Get(fmt.Sprintf("http://localhost%s", port))
				return err
			}).ShouldNot(HaveOccurred())

			Expect(gw.Stop()).To(BeTrue())

			Eventually(func() error {
				_, err := http.Get(fmt.Sprintf("http://localhost%s", port))
				return err
			}).Should(HaveOccurred())
		})

		It("does not stop twice", func() {
			Expect(gw.Stop()).To(BeTrue())
			Expect(gw.Stop()).To(BeFalse())
		})
	})

	Describe("AddService", func() {
		It("adds a new service", func() {
			Expect(gw.Services[spec.Host]).To(BeNil())

			Expect(gw.AddService(logger, spec)).To(Succeed())
			Expect(fakeReverseProxyCreator.CreateCallCount()).To(Equal(1))
			_, serviceSpec := fakeReverseProxyCreator.CreateArgsForCall(0)
			Expect(serviceSpec).To(Equal(spec))

			Expect(gw.Services[spec.Host]).NotTo(BeNil())
		})

		Context("when fails to create a service hostr", func() {
			BeforeEach(func() {
				fakeReverseProxyCreator.CreateReturns(nil, errors.New("failed to create hostr"))
			})

			It("returns an error", func() {
				Expect(gw.AddService(logger, spec)).To(MatchError(ContainSubstring("failed to create hostr")))
			})
		})
	})

	Describe("RemoveService", func() {
		BeforeEach(func() {
			Expect(gw.Services[spec.Host]).To(BeNil())
			Expect(gw.AddService(logger, spec)).To(Succeed())
			Expect(gw.Services[spec.Host]).NotTo(BeNil())
		})

		It("removes an existing service", func() {
			Expect(gw.RemoveService(logger, spec.Host)).To(Succeed())
			Expect(gw.Services[spec.Host]).To(BeNil())
		})

		Context("when service is not found", func() {
			It("returns an error", func() {
				Expect(gw.RemoveService(logger, "not-found")).To(HaveOccurred())
			})
		})
	})

	Describe("ServeHTTP", func() {
		var spec gateway.ReverseProxySpec

		BeforeEach(func() {
			spec = gateway.ReverseProxySpec{
				Host:     "my-host.apihub.dev",
				Backends: []string{},
			}

			fakeReverseProxyCreator.CreateReturns(fakeReverseProxy, nil)
			Expect(gw.AddService(logger, spec)).To(Succeed())
		})

		It("forwards the request to the reverse proxy", func() {
			rw := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodGet, "http://my-host.apihub.dev/ping", nil)
			Expect(err).NotTo(HaveOccurred())
			gw.ServeHTTP(rw, req)
			Expect(fakeReverseProxy.ServeHTTPCallCount()).To(Equal(1))
		})

		It("returns page not found when service is not found", func() {
			rw := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodGet, "http://not-found.example.com/ping", nil)
			Expect(err).NotTo(HaveOccurred())
			gw.ServeHTTP(rw, req)
			Expect(rw.Body.String()).To(MatchRegexp(`{"error":"not_found","error_description":"The requested resource could not be found but may be available again in the future."}`))
		})
	})
})
