package client_test

import (
	"github.com/apihub/apihub"
	"github.com/apihub/apihub/client"
	"github.com/apihub/apihub/client/connection/connectionfakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Service", func() {
	var (
		fakeConnection *connectionfakes.FakeConnection
		cli            apihub.Client
	)

	BeforeEach(func() {
		fakeConnection = new(connectionfakes.FakeConnection)
		cli = client.New(fakeConnection)
	})

	Describe("Handle", func() {
		It("returns service's handle", func() {
			spec := apihub.ServiceSpec{
				Handle: "my-handle",
			}

			fakeConnection.AddServiceReturns(spec, nil)
			service, err := cli.AddService(spec)
			Expect(err).NotTo(HaveOccurred())
			Expect(service.Handle()).To(Equal("my-handle"))
		})
	})

	Describe("Timeout", func() {
		It("returns service's timeout", func() {
			spec := apihub.ServiceSpec{
				Handle: "my-handle",
			}

			fakeConnection.AddServiceReturns(spec, nil)
			service, err := cli.AddService(spec)
			Expect(err).NotTo(HaveOccurred())
			Expect(service.Handle()).To(Equal("my-handle"))
		})
	})

})
