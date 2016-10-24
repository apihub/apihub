package client_test

import (
	"errors"

	"github.com/apihub/apihub"
	"github.com/apihub/apihub/client"
	"github.com/apihub/apihub/client/connection/connectionfakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Client", func() {

	var (
		cli            apihub.Client
		fakeConnection *connectionfakes.FakeConnection
	)

	BeforeEach(func() {
		fakeConnection = new(connectionfakes.FakeConnection)
		cli = client.New(fakeConnection)
	})

	Describe("AddService", func() {
		It("sends a request to add a service", func() {
			spec := apihub.ServiceSpec{
				Handle: "my-handle",
			}

			fakeConnection.AddServiceReturns(spec, nil)

			service, err := cli.AddService(spec)
			Expect(err).NotTo(HaveOccurred())
			Expect(service.Handle()).To(Equal("my-handle"))
			Expect(fakeConnection.AddServiceArgsForCall(0)).To(Equal(spec))
		})

		Context("when the request fails", func() {
			BeforeEach(func() {
				fakeConnection.AddServiceReturns(apihub.ServiceSpec{}, errors.New("failed to add service"))
			})

			It("returns an error", func() {
				spec := apihub.ServiceSpec{
					Handle: "my-handle",
				}
				_, err := cli.AddService(spec)
				Expect(err).To(HaveOccurred())
			})
		})

	})

})
