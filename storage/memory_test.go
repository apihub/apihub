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

		Context("when service is not found", func() {
			It("returns an error", func() {
				_, err := store.FindServiceByHandle("invalid-handle")
				Expect(err).To(MatchError(ContainSubstring("service not found")))
			})
		})
	})

	Describe("Services", func() {
		BeforeEach(func() {
			Expect(store.UpsertService(spec)).To(Succeed())
		})

		It("lists all services", func() {
			services, err := store.Services()
			Expect(err).NotTo(HaveOccurred())
			Expect(services).To(ConsistOf(spec))
		})
	})

	Describe("RemoveService", func() {
		BeforeEach(func() {
			Expect(store.UpsertService(spec)).To(Succeed())
		})

		It("removes service by handle", func() {
			Expect(store.RemoveService(spec.Handle)).NotTo(HaveOccurred())
		})

		Context("when service is not found", func() {
			It("returns an error", func() {
				Expect(store.RemoveService("invalid-handle")).To(MatchError(ContainSubstring("service not found")))
			})
		})
	})
})
