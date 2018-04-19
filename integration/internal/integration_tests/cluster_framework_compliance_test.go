package integration_tests

import (
	"io/ioutil"
	"path/filepath"

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

	It("Manages a configured etcd directory", func() {
		dir, err := ioutil.TempDir("", "")
		Expect(err).NotTo(HaveOccurred())

		dataDir := filepath.Join(dir, "etcd-test-dir")
		Expect(dataDir).NotTo(BeAnExistingFile())

		fixture = &integration.ControlPlane{}
		err = fixture.Setup(cluster.Config{Etcd: cluster.Etcd{DataDir: dataDir}})
		Expect(err).NotTo(HaveOccurred())

		Expect(dataDir).To(BeADirectory())

		Expect(fixture.TearDown()).To(Succeed())
	})
})
