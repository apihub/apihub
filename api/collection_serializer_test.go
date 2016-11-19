package api_test

import (
	"github.com/apihub/apihub"
	"github.com/apihub/apihub/api"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CollectionSerializer", func() {
	var (
		specs []*apihub.ServiceSpec
	)

	BeforeEach(func() {
		specs = []*apihub.ServiceSpec{
			&apihub.ServiceSpec{Host: "my-host.apihub.dev"},
		}
	})

	Describe("Collection", func() {
		It("returns a collection instance", func() {
			collection := api.Collection(specs, len(specs))
			Expect(collection.Items).To(Equal(specs))
			Expect(collection.Count).To(Equal(len(specs)))
		})
	})
})
