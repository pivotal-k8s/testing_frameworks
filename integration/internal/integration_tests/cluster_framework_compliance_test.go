package integration_tests

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"sigs.k8s.io/testing_frameworks/cluster"
	"sigs.k8s.io/testing_frameworks/integration"
)

var _ = Describe("Cluster Framework Compliance", func() {
	var fixture cluster.Fixture

	AfterEach(func() {
		Expect(fixture.TearDown()).To(Succeed())
	})

	It("Successfully manages the control plane lifecycle", func() {
		var err error

		fixture = &integration.ControlPlane{}
		By("Starting all the control plane processes")
		err = fixture.Setup(cluster.Config{})
		Expect(err).NotTo(HaveOccurred(), "Expected controlPlane to start successfully")

		apiURL := fixture.ClientConfig()
		isAPIServerListening := isSomethingListeningOnPort(apiURL.Host)
		Expect(isAPIServerListening()).To(BeTrue())
	})
})
