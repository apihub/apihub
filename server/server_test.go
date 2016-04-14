package server_test

import (
	"errors"
	"fmt"

	"github.com/apihub/apihub/apihubfakes"
	"github.com/apihub/apihub/server"
	"github.com/pivotal-golang/lager"
	"github.com/pivotal-golang/lager/lagertest"

	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/config"
	. "github.com/onsi/gomega"
)

var _ = Describe("The Apihub Server", func() {
	var (
		fakeBackend  *apihubfakes.FakeBackend
		listenAddr   string
		log          lager.Logger
		apihubServer *server.ApihubServer
	)

	BeforeEach(func() {
		listenAddr = fmt.Sprintf(":%d", 8080+config.GinkgoConfig.ParallelNode)
		log = lagertest.NewTestLogger("apihub-test")
	})

	JustBeforeEach(func() {
		fakeBackend = new(apihubfakes.FakeBackend)
		apihubServer = server.New(log, listenAddr, fakeBackend)
	})

	Describe("when starting up the server", func() {
		It("starts the backend", func() {
			err := apihubServer.Start()
			defer apihubServer.Stop()

			Expect(err).NotTo(HaveOccurred())
			Expect(fakeBackend.StartCallCount()).To(Equal(1))
		})

		Context("when fails to start the backend", func() {
			JustBeforeEach(func() {
				fakeBackend.StartReturns(errors.New("Boom!"))
			})

			It("returns an error", func() {
				err := apihubServer.Start()
				defer apihubServer.Stop()
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("when shutting down the server", func() {
		It("stops the backend", func() {
			err := apihubServer.Start()
			Expect(err).NotTo(HaveOccurred())
			apihubServer.Stop()
			Expect(fakeBackend.StopCallCount()).To(Equal(1))
		})

		It("stops the Apihub server", func() {
			err := apihubServer.Start()
			Expect(err).NotTo(HaveOccurred())
			apihubServer.Stop()
			Expect(fakeBackend.StopCallCount()).To(Equal(1))
			Eventually(func() error {
				for {
					_, err := apihubServer.Accept()
					return err
				}
			}).Should(HaveOccurred())
		})

		Context("when fails to stop the backend", func() {
			JustBeforeEach(func() {
				fakeBackend.StopReturns(errors.New("Boom!"))
			})

			It("returns an error", func() {
				err := apihubServer.Start()
				Expect(err).NotTo(HaveOccurred())

				err = apihubServer.Stop()
				Expect(err).To(HaveOccurred())
			})
		})
	})

})
