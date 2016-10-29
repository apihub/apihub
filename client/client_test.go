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
				Expect(err).To(MatchError(ContainSubstring("failed to add service")))
			})
		})
	})

	Describe("Services", func() {
		It("sends a request to list services", func() {
			fakeConnection.ServicesReturns([]apihub.ServiceSpec{
				apihub.ServiceSpec{
					Handle: "my-handle",
				},
			}, nil)

			services, err := cli.Services()
			Expect(err).NotTo(HaveOccurred())
			Expect(len(services)).To(Equal(1))
			Expect(services[0].Handle()).To(Equal("my-handle"))
		})

		Context("when the request fails", func() {
			BeforeEach(func() {
				fakeConnection.ServicesReturns(nil, errors.New("failed to list services"))
			})

			It("returns an error", func() {
				_, err := cli.Services()
				Expect(err).To(MatchError(ContainSubstring("failed to list services")))
			})
		})
	})

	Describe("RemoveService", func() {
		It("sends a request to remove a service", func() {
			fakeConnection.RemoveServiceReturns(nil)

			err := cli.RemoveService("my-handle")
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeConnection.RemoveServiceArgsForCall(0)).To(Equal("my-handle"))
		})

		Context("when the request fails", func() {
			BeforeEach(func() {
				fakeConnection.RemoveServiceReturns(errors.New("failed to remove service: `my-handle`"))
			})

			It("returns an error", func() {
				err := cli.RemoveService("my-handle")
				Expect(err).To(MatchError(ContainSubstring("failed to remove service: `my-handle`")))
			})
		})
	})

	Describe("FindService", func() {
		It("sends a request to find a service", func() {
			spec := apihub.ServiceSpec{
				Handle: "my-handle",
			}

			fakeConnection.FindServiceReturns(spec, nil)
			service, err := cli.FindService("my-handle")
			Expect(err).NotTo(HaveOccurred())
			Expect(service.Handle()).To(Equal("my-handle"))
			Expect(fakeConnection.FindServiceArgsForCall(0)).To(Equal("my-handle"))
		})

		Context("when the request fails", func() {
			BeforeEach(func() {
				fakeConnection.FindServiceReturns(apihub.ServiceSpec{}, errors.New("failed to find service: `my-handle`"))
			})

			It("returns an error", func() {
				_, err := cli.FindService("my-handle")
				Expect(err).To(MatchError(ContainSubstring("failed to find service: `my-handle`")))
			})
		})
	})
})
