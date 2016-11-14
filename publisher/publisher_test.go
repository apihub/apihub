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
		config apihub.ServiceConfig
	)

	BeforeEach(func() {
		logger = lagertest.NewTestLogger("publisher-test")

		config = apihub.ServiceConfig{
			ServiceSpec: apihub.ServiceSpec{
				Handle: "my-handle",
				Backends: []apihub.BackendInfo{
					apihub.BackendInfo{
						Address: "http://server-a",
					},
				},
			},
		}

		pub = publisher.NewPublisher(consulClient)
	})

	Describe("Publish", func() {
		It("publishes a service", func() {
			Expect(pub.Publish(logger, config)).To(Succeed())

			kvp, _, err := consulClient.KV().Get(config.ServiceSpec.Handle, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(kvp.Key).To(Equal("my-handle"))

			spec, err := json.Marshal(config.ServiceSpec)
			Expect(err).NotTo(HaveOccurred())
			Expect(kvp.Value).To(Equal(spec))
		})
	})
})
