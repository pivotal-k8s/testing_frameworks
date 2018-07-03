package cluster_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "sigs.k8s.io/testing_frameworks/cluster"
)

var _ = Describe("Cluster", func() {
	It("Defaults Etcd.Local", func() {
		c := Config{}
		c = DoDefaulting(c)
		Expect(c.Etcd.Local).NotTo(BeNil())
	})
})
