package integration_test

import (
	"fmt"

	"code.cloudfoundry.org/lager/lagertest"

	"github.com/apihub/apihub"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Service", func() {
	var (
		client      *RunningApihub
		addressAPI  string
		portGateway string
		spec        apihub.ServiceSpec
		logger      *lagertest.TestLogger
	)

	BeforeEach(func() {
		addressAPI = fmt.Sprintf("/tmp/apihub_api_%d.sock",
			GinkgoParallelNode())
		portGateway = fmt.Sprintf(":%d", 9000+GinkgoParallelNode())
		logger = lagertest.NewTestLogger("services-test")

		spec = apihub.ServiceSpec{
			Handle:   "my-service",
			Disabled: true,
			Timeout:  10,
			Backends: []apihub.BackendInfo{
				apihub.BackendInfo{
					Address:          "http://server-a",
					HeartBeatAddress: "http://server-a/healthcheck",
					HeartBeatTimeout: 3,
				},
			},
		}
	})

	JustBeforeEach(func() {
		client = startApihub("unix", addressAPI, portGateway)
	})

	AfterEach(func() {
		Expect(client.Stop()).To(Succeed())
	})

	Describe("AddService", func() {
		It("adds a new service", func() {
			service, err := client.AddService(spec)
			Expect(err).NotTo(HaveOccurred())
			Expect(service.Handle()).To(Equal("my-service"))
		})

		Context("when there's another service for given handle", func() {
			JustBeforeEach(func() {
				service, err := client.AddService(spec)
				Expect(err).NotTo(HaveOccurred())
				Expect(service.Handle()).To(Equal("my-service"))
			})

			It("returns an error message with bad request", func() {
				_, err := client.AddService(spec)
				Expect(err).To(MatchError(ContainSubstring("handle already in use")))
			})
		})
	})

	Describe("Services", func() {
		JustBeforeEach(func() {
			service, err := client.AddService(spec)
			Expect(err).NotTo(HaveOccurred())
			Expect(service.Handle()).To(Equal("my-service"))
		})

		It("lists services", func() {
			services, err := client.Services()
			Expect(err).NotTo(HaveOccurred())
			Expect(len(services)).To(Equal(1))
			Expect(services[0].Handle()).To(Equal("my-service"))
		})
	})

	Describe("RemoveService", func() {
		JustBeforeEach(func() {
			_, err := client.AddService(spec)
			Expect(err).NotTo(HaveOccurred())
		})

		It("removes a service", func() {
			err := client.RemoveService("my-service")
			Expect(err).NotTo(HaveOccurred())

			services, err := client.Services()
			Expect(err).NotTo(HaveOccurred())
			Expect(len(services)).To(Equal(0))
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
			service, err := client.FindService("my-service")
			Expect(err).NotTo(HaveOccurred())
			Expect(service.Handle()).To(Equal("my-service"))
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

			service, err := client.UpdateService("my-service", spec)
			Expect(err).NotTo(HaveOccurred())

			service, err = client.FindService("my-service")
			Expect(err).NotTo(HaveOccurred())
			backends, err := service.Backends()
			Expect(err).NotTo(HaveOccurred())
			Expect(backends[0].Address).To(Equal("http://server-b"))
		})

		Context("when service is not found", func() {
			It("returns an error", func() {
				_, err := client.UpdateService("invalid-handle", spec)
				Expect(err).To(MatchError(ContainSubstring("Failed to find service.")))
			})
		})
	})
})
