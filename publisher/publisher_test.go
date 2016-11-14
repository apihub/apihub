package publisher_test

import (
	"encoding/json"

	"code.cloudfoundry.org/lager/lagertest"
	"github.com/apihub/apihub"
	"github.com/apihub/apihub/publisher"

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
			Handle: "my-handle",
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
			Expect(pub.Publish(logger, spec)).To(Succeed())

			kvp, _, err := consulClient.KV().Get(spec.Handle, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(kvp.Key).To(Equal("my-handle"))

			spec, err := json.Marshal(spec)
			Expect(err).NotTo(HaveOccurred())
			Expect(kvp.Value).To(Equal(spec))
		})
	})
})
