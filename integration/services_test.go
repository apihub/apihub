package integration_test

import (
	"fmt"

	"github.com/apihub/apihub"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Service", func() {
	var (
		client  *RunningApihub
		address string
		spec    apihub.ServiceSpec
	)

	BeforeEach(func() {
		address = fmt.Sprintf("/tmp/apihub_%d.sock",
			GinkgoParallelNode())

		spec = apihub.ServiceSpec{
			Handle:   "my-service",
			Disabled: true,
			Timeout:  10,
			Backends: []apihub.BackendInfo{
				apihub.BackendInfo{
					Name:             "server-a",
					Address:          "http://server-a",
					HeartBeatAddress: "http://server-a/healthcheck",
					HeartBeatTimeout: 3,
					HeartBeatRetry:   2,
				},
			},
		}
	})

	JustBeforeEach(func() {
		client = startApihub("unix", address)
	})

	AfterEach(func() {
		Expect(client.Stop()).To(Succeed())
	})

	Describe("Add a service", func() {
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
				Expect(err).To(MatchError(ContainSubstring("Handle already in use.")))
			})
		})
	})

	Describe("List existing services", func() {
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

	Describe("Remove a service", func() {
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
			It("returns an error message with bad request", func() {
				err := client.RemoveService("invalid-handle")
				Expect(err).To(MatchError(ContainSubstring("Handle not found.")))
			})
		})
	})
})
