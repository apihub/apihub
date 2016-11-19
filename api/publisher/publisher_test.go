package publisher_test

import (
	"encoding/json"
	"fmt"

	"code.cloudfoundry.org/lager/lagertest"
	"github.com/apihub/apihub"
	"github.com/apihub/apihub/api/publisher"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Publisher", func() {

	var (
		logger *lagertest.TestLogger
		pub    *publisher.Publisher
		spec   apihub.ServiceSpec
	)

	BeforeEach(func() {
		logger = lagertest.NewTestLogger("publisher-test")

		spec = apihub.ServiceSpec{
			Host: "my-host.apihub.dev",
			Backends: []apihub.BackendInfo{
				apihub.BackendInfo{
					Address: "http://server-a",
				},
			},
		}

		pub = publisher.NewPublisher(consulClient)
	})

	Describe("Publish", func() {
		It("publishes a service", func() {
			Expect(pub.Publish(logger, apihub.SERVICES_PREFIX, spec)).To(Succeed())

			key := fmt.Sprintf("%s%s", apihub.SERVICES_PREFIX, spec.Host)
			kvp, _, err := consulClient.KV().Get(key, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(kvp).NotTo(BeNil())

			spec, err := json.Marshal(spec)
			Expect(err).NotTo(HaveOccurred())
			Expect(kvp.Value).To(Equal(spec))
		})
	})

	Describe("Unpublish", func() {
		It("unpublishes a service", func() {
			Expect(pub.Publish(logger, apihub.SERVICES_PREFIX, spec)).To(Succeed())
			Expect(pub.Unpublish(logger, apihub.SERVICES_PREFIX, spec.Host)).To(Succeed())

			key := fmt.Sprintf("%s%s", apihub.SERVICES_PREFIX, spec.Host)
			kvp, _, err := consulClient.KV().Get(key, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(kvp).To(BeNil())
		})

		Context("when service is not found", func() {
			It("does not return an error", func() {
				Expect(pub.Unpublish(logger, apihub.SERVICES_PREFIX, "not-found-key")).To(Succeed())
			})
		})
	})
})
