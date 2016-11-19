package subscriber_test

import (
	"code.cloudfoundry.org/lager/lagertest"
	"github.com/apihub/apihub"
	"github.com/apihub/apihub/api/publisher"
	"github.com/apihub/apihub/gateway/subscriber"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Subscriber", func() {
	var (
		logger     *lagertest.TestLogger
		pub        *publisher.Publisher
		sub        *subscriber.Subscriber
		spec       apihub.ServiceSpec
		servicesCh chan apihub.ServiceSpec
		stop       chan struct{}
	)

	BeforeEach(func() {
		logger = lagertest.NewTestLogger("subscriber-test")

		spec = apihub.ServiceSpec{
			Handle: "my-handle",
			Backends: []apihub.BackendInfo{
				apihub.BackendInfo{
					Address: "http://server-a",
				},
			},
		}

		pub = publisher.NewPublisher(consulClient)
		sub = subscriber.NewSubscriber(consulClient)
		servicesCh = make(chan apihub.ServiceSpec)
		stop = make(chan struct{})
	})

	AfterEach(func() {
		Expect(consulRunner.Reset()).To(Succeed())
	})

	Describe("Subscribe", func() {
		It("loads services already published", func() {
			Expect(pub.Publish(logger, apihub.SERVICES_PREFIX, spec)).To(Succeed())

			go func() {
				err := sub.Subscribe(logger, apihub.SERVICES_PREFIX, servicesCh, stop)
				Expect(err).NotTo(HaveOccurred())
			}()

			Eventually(servicesCh).Should(Receive(Equal(spec)))
			Consistently(stop).ShouldNot(BeClosed())
			close(stop)
			Eventually(stop).Should(BeClosed())
		})

		It("receives new services", func() {
			go func() {
				err := sub.Subscribe(logger, apihub.SERVICES_PREFIX, servicesCh, stop)
				Expect(err).NotTo(HaveOccurred())
			}()

			anotherSpec := apihub.ServiceSpec{
				Handle: "my-second-handle",
			}
			Expect(pub.Publish(logger, apihub.SERVICES_PREFIX, anotherSpec)).To(Succeed())
			Eventually(servicesCh).Should(Receive(Equal(anotherSpec)))
			Consistently(stop).ShouldNot(BeClosed())
			close(stop)
			Eventually(stop).Should(BeClosed())
		})

		It("updates existing services", func() {
			go func() {
				err := sub.Subscribe(logger, apihub.SERVICES_PREFIX, servicesCh, stop)
				Expect(err).NotTo(HaveOccurred())
			}()
			Expect(pub.Publish(logger, apihub.SERVICES_PREFIX, spec)).To(Succeed())
			Eventually(servicesCh).Should(Receive(Equal(spec)))

			spec.Handle = "another-handle"
			Expect(pub.Publish(logger, apihub.SERVICES_PREFIX, spec)).To(Succeed())

			service := <-servicesCh
			Expect(service.Handle).To(Equal("another-handle"))

			Consistently(stop).ShouldNot(BeClosed())
			close(stop)
			Eventually(stop).Should(BeClosed())
		})

		Context("when the subscription is stopped", func() {
			It("closes services channel", func() {
				close(stop)
				Expect(sub.Subscribe(logger, apihub.SERVICES_PREFIX, servicesCh, stop)).To(Succeed())
				Eventually(servicesCh).Should(BeClosed())
			})
		})

		Context("when an error occurs", func() {
			It("retries to connect", func() {
				go func() {
					err := sub.Subscribe(logger, apihub.SERVICES_PREFIX, servicesCh, stop)
					Expect(err).NotTo(HaveOccurred())
				}()

				consulRunner.Stop()
				Consistently(servicesCh).ShouldNot(Receive())

				consulRunner.Start()
				consulRunner.WaitUntilReady()

				spec := apihub.ServiceSpec{
					Handle: "my-retry",
				}
				Expect(pub.Publish(logger, apihub.SERVICES_PREFIX, spec)).To(Succeed())
				Eventually(servicesCh).Should(Receive(Equal(spec)))
			})
		})
	})
})
