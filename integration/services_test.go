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

})
