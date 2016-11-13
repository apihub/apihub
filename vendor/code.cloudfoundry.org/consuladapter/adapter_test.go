package consuladapter_test

import (
	"code.cloudfoundry.org/consuladapter"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Adapter", func() {
	Describe("Parse", func() {
		It("errors when passed an invalid URL", func() {
			_, _, err := consuladapter.Parse(":/")
			Expect(err).To(HaveOccurred())
		})

		It("errors when passed a scheme that is not http or https", func() {
			_, _, err := consuladapter.Parse("sftp://")
			Expect(err).To(HaveOccurred())
		})

		It("errors when passed an empty host", func() {
			_, _, err := consuladapter.Parse("http:///")
			Expect(err).To(HaveOccurred())
		})

		It("returns the scheme and address", func() {
			scheme, address, err := consuladapter.Parse("https://1.2.3.4:5678")
			Expect(err).NotTo(HaveOccurred())
			Expect(scheme).To(Equal("https"))
			Expect(address).To(Equal("1.2.3.4:5678"))
		})
	})
})
