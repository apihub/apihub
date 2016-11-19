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
				Host: "my-host.apihub.dev",
			}

			fakeConnection.AddServiceReturns(spec, nil)

			service, err := cli.AddService(spec)
			Expect(err).NotTo(HaveOccurred())
			Expect(service.Host()).To(Equal("my-host.apihub.dev"))
			Expect(fakeConnection.AddServiceArgsForCall(0)).To(Equal(spec))
		})

		Context("when the request fails", func() {
			BeforeEach(func() {
				fakeConnection.AddServiceReturns(apihub.ServiceSpec{}, errors.New("failed to add service"))
			})

			It("returns an error", func() {
				spec := apihub.ServiceSpec{
					Host: "my-host.apihub.dev",
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
					Host: "my-host.apihub.dev",
				},
			}, nil)

			services, err := cli.Services()
			Expect(err).NotTo(HaveOccurred())
			Expect(len(services)).To(Equal(1))
			Expect(services[0].Host()).To(Equal("my-host.apihub.dev"))
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

			err := cli.RemoveService("my-host.apihub.dev")
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeConnection.RemoveServiceArgsForCall(0)).To(Equal("my-host.apihub.dev"))
		})

		Context("when the request fails", func() {
			BeforeEach(func() {
				fakeConnection.RemoveServiceReturns(errors.New("failed to remove service: `my-host.apihub.dev`"))
			})

			It("returns an error", func() {
				err := cli.RemoveService("my-host.apihub.dev")
				Expect(err).To(MatchError(ContainSubstring("failed to remove service: `my-host.apihub.dev`")))
			})
		})
	})

	Describe("FindService", func() {
		It("sends a request to find a service", func() {
			spec := apihub.ServiceSpec{
				Host: "my-host.apihub.dev",
			}

			fakeConnection.FindServiceReturns(spec, nil)
			service, err := cli.FindService("my-host.apihub.dev")
			Expect(err).NotTo(HaveOccurred())
			Expect(service.Host()).To(Equal("my-host.apihub.dev"))
			Expect(fakeConnection.FindServiceArgsForCall(0)).To(Equal("my-host.apihub.dev"))
		})

		Context("when the request fails", func() {
			BeforeEach(func() {
				fakeConnection.FindServiceReturns(apihub.ServiceSpec{}, errors.New("failed to find service: `my-host.apihub.dev`"))
			})

			It("returns an error", func() {
				_, err := cli.FindService("my-host.apihub.dev")
				Expect(err).To(MatchError(ContainSubstring("failed to find service: `my-host.apihub.dev`")))
			})
		})
	})

	Describe("UpdateService", func() {
		var spec apihub.ServiceSpec

		BeforeEach(func() {
			spec = apihub.ServiceSpec{
				Host: "my-host.apihub.dev",
			}
			fakeConnection.FindServiceReturns(spec, nil)
		})

		It("sends a request to updaate a service", func() {
			fakeConnection.UpdateServiceReturns(spec, nil)

			_, err := cli.UpdateService("my-host.apihub.dev", spec)
			Expect(err).NotTo(HaveOccurred())

			host, s := fakeConnection.UpdateServiceArgsForCall(0)
			Expect(host).To(Equal("my-host.apihub.dev"))
			Expect(s).To(Equal(spec))
		})

		Context("when the request fails", func() {
			BeforeEach(func() {
				fakeConnection.FindServiceReturns(apihub.ServiceSpec{}, errors.New("failed to find service: `my-host.apihub.dev`"))
			})

			It("returns an error", func() {
				_, err := cli.FindService("my-host.apihub.dev")
				Expect(err).To(MatchError(ContainSubstring("failed to find service: `my-host.apihub.dev`")))
			})
		})
	})
})
