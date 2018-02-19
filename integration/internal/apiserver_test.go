package internal_test

import (
	"net/url"

	. "github.com/kubernetes-sig-testing/frameworks/integration/internal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("GetAPIServerStartMessage()", func() {
	Context("when using a non tls URL", func() {
		It("generates valid start message", func() {
			url := url.URL{
				Scheme: "http",
				Host:   "some.insecure.apiserver:1234",
			}
			message := GetAPIServerStartMessage(url)
			Expect(message).To(Equal("Serving insecurely on some.insecure.apiserver:1234"))
		})
	})
	Context("when using a tls URL", func() {
		It("generates valid start message", func() {
			url := url.URL{
				Scheme: "https",
				Host:   "some.secure.apiserver:8443",
			}
			message := GetAPIServerStartMessage(url)
			Expect(message).To(Equal("Serving securely on some.secure.apiserver:8443"))
		})
	})
})
