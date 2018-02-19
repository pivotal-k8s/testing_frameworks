package integration

import (
	"net/url"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/kubernetes-sig-testing/frameworks/integration/internal"
)

var _ = Describe("processStateToConnectionConfig", func() {
	It("generates a proper RemoteConnectionConfig", func() {
		ps := &internal.ProcessState{}
		ps.URL = url.URL{Scheme: "http", Host: "some.host.tld"}

		c, err := processStateToConnectionConfig(ps)
		Expect(err).NotTo(HaveOccurred())
		Expect(c.URL.String()).To(Equal("http://some.host.tld"))
	})

	Context("when process has not defaulted the URL yet", func() {
		It("propagates the error", func() {
			ps := &internal.ProcessState{}

			_, err := processStateToConnectionConfig(ps)
			Expect(err).To(MatchError(ContainSubstring("not bound to an URL yet")))
		})
	})

	Context("when there is no process state", func() {
		It("returns a error", func() {
			_, err := processStateToConnectionConfig(nil)
			Expect(err).To(MatchError(ContainSubstring("no process state")))
		})
	})
})
