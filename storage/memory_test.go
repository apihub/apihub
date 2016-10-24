package storage_test

import (
	"github.com/apihub/apihub"
	"github.com/apihub/apihub/storage"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Memory", func() {
	var (
		store apihub.Storage
		spec  apihub.ServiceSpec
	)

	BeforeEach(func() {
		store = storage.New()
		spec = apihub.ServiceSpec{Handle: "my-handle"}
	})

	Describe("UpsertService", func() {
		It("adds a service", func() {
			Expect(store.UpsertService(spec)).To(Succeed())
		})
	})

	Describe("FindServiceByHandle", func() {
		BeforeEach(func() {
			Expect(store.UpsertService(spec)).To(Succeed())
		})

		It("finds a service", func() {
			found, err := store.FindServiceByHandle("my-handle")
			Expect(err).NotTo(HaveOccurred())
			Expect(found).To(Equal(spec))
		})
	})
})
