package integration_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"code.cloudfoundry.org/lager/lagertest"

	"github.com/apihub/apihub"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Service", func() {
	var (
		client      *RunningApihub
		addressAPI  string
		portGateway int
		spec        apihub.ServiceSpec
		logger      *lagertest.TestLogger
		testServer  *httptest.Server
	)

	BeforeEach(func() {
		addressAPI = fmt.Sprintf("/tmp/apihub_api_%d.sock",
			GinkgoParallelNode())
		portGateway = 9000 + GinkgoParallelNode()
		logger = lagertest.NewTestLogger("services-test")

		testServer = httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			rw.Write([]byte("Hello World!"))
		}))

		spec = apihub.ServiceSpec{
			Handle: fmt.Sprintf("my-service-%d", GinkgoParallelNode()),
			Backends: []apihub.BackendInfo{
				apihub.BackendInfo{
					Address: "http://" + testServer.Listener.Addr().String(),
				},
			},
		}
	})

	JustBeforeEach(func() {
		client = startApihub("unix", addressAPI, portGateway)
	})

	AfterEach(func() {
		testServer.Close()
		Expect(client.Stop()).To(Succeed())
	})

	fireRequest := func(portGateway int, handle string) *http.Response {
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://127.0.0.1:%d", portGateway), nil)
		Expect(err).NotTo(HaveOccurred())
		req.Host = fmt.Sprintf("%s.apihub.dev", handle)

		c := &http.Client{}
		resp, err := c.Do(req)
		Expect(err).NotTo(HaveOccurred())

		return resp
	}

	Describe("AddService", func() {
		It("adds a new service", func() {
			service, err := client.AddService(spec)
			Expect(err).NotTo(HaveOccurred())
			Expect(service.Handle()).To(Equal(spec.Handle))
		})

		It("proxies the request to the service endpoint", func() {
			service, err := client.AddService(spec)
			Expect(err).NotTo(HaveOccurred())
			resp := fireRequest(portGateway, service.Handle())
			body, err := ioutil.ReadAll(resp.Body)
			Expect(err).NotTo(HaveOccurred())
			Expect(string(body)).To(Equal("Hello World!"))
		})

		Context("when there's another service for given handle", func() {
			JustBeforeEach(func() {
				_, err := client.AddService(spec)
				Expect(err).NotTo(HaveOccurred())
			})

			It("returns an error message with bad request", func() {
				_, err := client.AddService(spec)
				Expect(err).To(MatchError(ContainSubstring("handle already in use")))
			})
		})
	})

	Describe("Services", func() {
		JustBeforeEach(func() {
			_, err := client.AddService(spec)
			Expect(err).NotTo(HaveOccurred())
		})

		It("lists services", func() {
			services, err := client.Services()
			Expect(err).NotTo(HaveOccurred())
			Expect(len(services)).To(Equal(1))
			Expect(services[0].Handle()).To(Equal(spec.Handle))
		})
	})

	Describe("RemoveService", func() {
		JustBeforeEach(func() {
			_, err := client.AddService(spec)
			Expect(err).NotTo(HaveOccurred())
		})

		It("removes a service", func() {
			err := client.RemoveService(spec.Handle)
			Expect(err).NotTo(HaveOccurred())

			services, err := client.Services()
			Expect(err).NotTo(HaveOccurred())
			Expect(len(services)).To(Equal(0))
		})

		It("unpublishes the service", func() {
			// Check if service is up and running
			resp := fireRequest(portGateway, spec.Handle)
			Eventually(resp.StatusCode).Should(Equal(http.StatusOK))
			body, err := ioutil.ReadAll(resp.Body)
			Expect(err).NotTo(HaveOccurred())
			Expect(string(body)).To(Equal("Hello World!"))

			// Remove service
			err = client.RemoveService(spec.Handle)
			Expect(err).NotTo(HaveOccurred())

			resp = fireRequest(portGateway, spec.Handle)
			Eventually(resp.StatusCode).Should(Equal(http.StatusNotFound))
		})

		Context("when service is not found", func() {
			It("returns an error", func() {
				err := client.RemoveService("invalid-handle")
				Expect(err).To(MatchError(ContainSubstring("Handle not found.")))
			})
		})
	})

	Describe("FindService", func() {
		JustBeforeEach(func() {
			_, err := client.AddService(spec)
			Expect(err).NotTo(HaveOccurred())
		})

		It("finds a service by handle", func() {
			service, err := client.FindService(spec.Handle)
			Expect(err).NotTo(HaveOccurred())
			Expect(service.Handle()).To(Equal(spec.Handle))
		})

		Context("when service is not found", func() {
			It("returns an error", func() {
				_, err := client.FindService("invalid-handle")
				Expect(err).To(MatchError(ContainSubstring("Failed to find service.")))
			})
		})
	})

	Describe("UpdateService", func() {
		JustBeforeEach(func() {
			_, err := client.AddService(spec)
			Expect(err).NotTo(HaveOccurred())
		})

		It("update an existing service by handle", func() {
			spec.Backends = []apihub.BackendInfo{
				apihub.BackendInfo{
					Address:          "http://server-b",
					HeartBeatAddress: "http://server-b/healthcheck",
					HeartBeatTimeout: 3,
				},
			}

			service, err := client.UpdateService(spec.Handle, spec)
			Expect(err).NotTo(HaveOccurred())

			service, err = client.FindService(spec.Handle)
			Expect(err).NotTo(HaveOccurred())
			backends, err := service.Backends()
			Expect(err).NotTo(HaveOccurred())
			Expect(backends[0].Address).To(Equal("http://server-b"))
		})

		Context("when the service is enabled", func() {
			BeforeEach(func() {
				spec.Disabled = true
			})

			It("proxies the request to the service endpoint", func() {
				resp := fireRequest(portGateway, spec.Handle)
				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))

				spec.Disabled = false
				_, err := client.UpdateService(spec.Handle, spec)
				Expect(err).NotTo(HaveOccurred())

				resp = fireRequest(portGateway, spec.Handle)
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				body, err := ioutil.ReadAll(resp.Body)
				Expect(err).NotTo(HaveOccurred())
				Expect(string(body)).To(Equal("Hello World!"))
			})
		})

		Context("when the service is disabled", func() {
			It("stops proxing the request to the service endpoint", func() {
				// Check if service is up and running
				resp := fireRequest(portGateway, spec.Handle)
				Eventually(resp.StatusCode).Should(Equal(http.StatusOK))
				body, err := ioutil.ReadAll(resp.Body)
				Expect(err).NotTo(HaveOccurred())
				Expect(string(body)).To(Equal("Hello World!"))

				// Disable service
				spec.Disabled = true
				_, err = client.UpdateService(spec.Handle, spec)
				Expect(err).NotTo(HaveOccurred())

				resp = fireRequest(portGateway, spec.Handle)
				Eventually(resp.StatusCode).Should(Equal(http.StatusNotFound))
			})
		})

		Context("when service is not found", func() {
			It("returns an error", func() {
				_, err := client.UpdateService("invalid-handle", spec)
				Expect(err).To(MatchError(ContainSubstring("Failed to find service.")))
			})
		})
	})
})
