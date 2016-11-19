package client_test

import (
	"errors"
	"time"

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
		service        apihub.Service
		spec           apihub.ServiceSpec
	)

	BeforeEach(func() {
		fakeConnection = new(connectionfakes.FakeConnection)
		cli = client.New(fakeConnection)

		spec = apihub.ServiceSpec{
			Handle:   "my-handle",
			Disabled: true,
			Timeout:  10,
			Backends: []apihub.BackendInfo{
				apihub.BackendInfo{
					Address: "http://server-a",
				},
				apihub.BackendInfo{
					Address: "http://server-b",
				},
			},
		}
	})

	JustBeforeEach(func() {
		var err error
		fakeConnection.AddServiceReturns(spec, nil)

		service, err = cli.AddService(spec)
		Expect(err).NotTo(HaveOccurred())
	})

	Describe("Handle", func() {
		It("returns service's handle", func() {
			Expect(service.Handle()).To(Equal("my-handle"))
		})
	})

	Describe("Info", func() {
		BeforeEach(func() {
			fakeConnection.FindServiceReturns(spec, nil)
		})

		It("returns service's info", func() {
			info, err := service.Info()
			Expect(err).NotTo(HaveOccurred())
			Expect(info.Handle).To(Equal(spec.Handle))
			Expect(info.Disabled).To(Equal(spec.Disabled))
			Expect(info.Timeout).To(Equal(spec.Timeout))
			Expect(info.Backends).To(ConsistOf(spec.Backends))
		})

		Context("when fails to get the info", func() {
			BeforeEach(func() {
				fakeConnection.FindServiceReturns(apihub.ServiceSpec{}, errors.New("fail to get the info"))
			})

			It("returns an error", func() {
				_, err := service.Info()
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("Backends", func() {
		BeforeEach(func() {
			fakeConnection.FindServiceReturns(spec, nil)
		})

		It("returns service's backends", func() {
			backends, err := service.Backends()
			Expect(err).NotTo(HaveOccurred())
			Expect(len(backends)).To(Equal(2))
		})

		Context("when fails to get the backend list", func() {
			BeforeEach(func() {
				fakeConnection.FindServiceReturns(apihub.ServiceSpec{}, errors.New("fail to get the backend list"))
			})

			It("returns an error", func() {
				_, err := service.Backends()
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("Start", func() {
		It("enables the service to start receiving requests", func() {
			Expect(service.Start()).To(Succeed())
			Expect(fakeConnection.UpdateServiceCallCount()).To(Equal(1))
			handle, serviceSpec := fakeConnection.UpdateServiceArgsForCall(0)
			Expect(handle).To(Equal(spec.Handle))
			Expect(spec.Disabled).To(BeTrue())
			Expect(serviceSpec.Disabled).To(BeFalse())
		})

		Context("when fails to enable", func() {
			BeforeEach(func() {
				fakeConnection.UpdateServiceReturns(apihub.ServiceSpec{}, errors.New("failed to update"))
			})

			It("returns an error", func() {
				Expect(service.Start()).To(MatchError(ContainSubstring("failed to update")))
			})
		})
	})

	Describe("Stop", func() {
		BeforeEach(func() {
			spec.Disabled = false
		})

		It("disables the service to stop receiving requests", func() {
			Expect(service.Stop()).To(Succeed())
			Expect(fakeConnection.UpdateServiceCallCount()).To(Equal(1))
			handle, serviceSpec := fakeConnection.UpdateServiceArgsForCall(0)
			Expect(handle).To(Equal(spec.Handle))
			Expect(spec.Disabled).To(BeFalse())
			Expect(serviceSpec.Disabled).To(BeTrue())
		})

		Context("when fails to disable", func() {
			BeforeEach(func() {
				fakeConnection.UpdateServiceReturns(apihub.ServiceSpec{}, errors.New("failed to update"))
			})

			It("returns an error", func() {
				Expect(service.Stop()).To(MatchError(ContainSubstring("failed to update")))
			})
		})
	})

	Describe("SetTimeout", func() {
		It("updates de timeout", func() {
			duration := time.Second * 10
			Expect(service.SetTimeout(duration)).To(Succeed())
			Expect(fakeConnection.UpdateServiceCallCount()).To(Equal(1))
			handle, serviceSpec := fakeConnection.UpdateServiceArgsForCall(0)
			Expect(handle).To(Equal(spec.Handle))
			Expect(serviceSpec.Timeout).To(Equal(duration))
		})

		Context("when fails to disable", func() {
			BeforeEach(func() {
				fakeConnection.UpdateServiceReturns(apihub.ServiceSpec{}, errors.New("failed to update"))
			})

			It("returns an error", func() {
				Expect(service.SetTimeout(time.Second)).To(MatchError(ContainSubstring("failed to update")))
			})
		})
	})
})
