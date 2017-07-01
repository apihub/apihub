package consuladapter_test

import (
	"code.cloudfoundry.org/consuladapter"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ConsulAdapter Client", func() {
	Context("with a strict HTTPS enabled consul cluster", func() {
		Context("with a HTTPS client", func() {
			Context("when valid values for certificates are provided", func() {
				It("is able to query consul", func() {
					consulClient = consulRunner.NewClient()
					_, err = consulClient.Status().Leader()
					Expect(err).ToNot(HaveOccurred())
				})
			})

			Context("when invalid values for certificates are provided", func() {
				It("is not able to query consul", func() {
					consulClient, err = consuladapter.NewTLSClientFromUrl(consulRunner.URL(), "", "", "")
					Expect(err).ToNot(HaveOccurred())
					_, err = consulClient.Status().Leader()
					Expect(err).To(HaveOccurred())
				})
			})

			Context("when an incorrect key is supplied", func() {
				It("is not able to create a client", func() {
					consulClient, err = consuladapter.NewTLSClientFromUrl(
						consulRunner.URL(),
						consulCACert,
						consulClientCert,
						invalidConsulClientKey,
					)
					Expect(err).To(HaveOccurred())
				})
			})
		})

		Context("with a HTTP client", func() {
			It("is not able to query consul on the HTTPS endpoint", func() {
				consulClient, err = consuladapter.NewClientFromUrl("http://" + consulRunner.Address())
				Expect(err).ToNot(HaveOccurred())
				_, err = consulClient.Status().Leader()
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
