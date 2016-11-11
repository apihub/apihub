package gateway_test

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"

	"code.cloudfoundry.org/lager/lagertest"
	"github.com/apihub/apihub/apihubfakes"
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
		fakeSubscriber          *apihubfakes.FakeServiceSubscriber
		fakeReverseProxyCreator *gatewayfakes.FakeReverseProxyCreator
		fakeReverseProxy        *gatewayfakes.FakeReverseProxy
	)

	BeforeEach(func() {
		logger = lagertest.NewTestLogger("test-gateway")
		port = fmt.Sprintf(":909%d", GinkgoParallelNode())
		fakeSubscriber = new(apihubfakes.FakeServiceSubscriber)
		fakeReverseProxyCreator = new(gatewayfakes.FakeReverseProxyCreator)
		fakeReverseProxy = new(gatewayfakes.FakeReverseProxy)

		gw = gateway.New(port, fakeSubscriber, fakeReverseProxyCreator)
	})

	var spec gateway.ReverseProxySpec

	BeforeEach(func() {
		spec = gateway.ReverseProxySpec{
			Handle:   "my-handle",
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
			_, err := http.Get(fmt.Sprintf("http://localhost%s", port))
			Expect(err).NotTo(HaveOccurred())

			Expect(gw.Stop()).To(BeTrue())

			_, err = http.Get(fmt.Sprintf("http://localhost%s", port))
			Expect(err).To(HaveOccurred())
		})

		It("does not stop twice", func() {
			Expect(gw.Stop()).To(BeTrue())
			Expect(gw.Stop()).To(BeFalse())
		})
	})

	Describe("AddService", func() {

		It("adds a new service", func() {
			Expect(gw.Services[spec.Handle]).To(BeNil())

			Expect(gw.AddService(logger, spec)).To(Succeed())
			Expect(fakeReverseProxyCreator.CreateCallCount()).To(Equal(1))
			_, serviceSpec := fakeReverseProxyCreator.CreateArgsForCall(0)
			Expect(serviceSpec).To(Equal(spec))

			Expect(gw.Services[spec.Handle]).NotTo(BeNil())
		})

		Context("when fails to create a service handler", func() {
			BeforeEach(func() {
				fakeReverseProxyCreator.CreateReturns(nil, errors.New("failed to create handler"))
			})

			It("returns an error", func() {
				Expect(gw.AddService(logger, spec)).To(MatchError(ContainSubstring("failed to create handler")))
			})
		})
	})

	Describe("ServeHTTP", func() {
		var spec gateway.ReverseProxySpec

		BeforeEach(func() {
			spec = gateway.ReverseProxySpec{
				Handle:   "my-handle",
				Backends: []string{},
			}

			fakeReverseProxyCreator.CreateReturns(fakeReverseProxy, nil)
			Expect(gw.AddService(logger, spec)).To(Succeed())
		})

		It("forwards the request to the reverse proxy", func() {
			rw := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodGet, "http://my-handle.example.com/ping", nil)
			Expect(err).NotTo(HaveOccurred())
			gw.ServeHTTP(rw, req)
			Expect(fakeReverseProxy.ServeHTTPCallCount()).To(Equal(1))
		})
	})
})
