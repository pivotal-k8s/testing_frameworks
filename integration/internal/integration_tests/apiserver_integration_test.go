package integration_tests

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/kubernetes-sig-testing/frameworks/integration"
)

var _ = Describe("APIServer", func() {
	Context("when no EtcdURL is provided", func() {
		It("does not panic", func() {
			apiServer := &APIServer{}

			starter := func() {
				Expect(apiServer.Start(RemoteConnectionConfig{})).To(
					MatchError(ContainSubstring("expected Etcd URL to be configured")),
				)
			}

			Expect(starter).NotTo(Panic())
		})
	})
})
