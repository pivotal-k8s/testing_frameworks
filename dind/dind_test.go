package dind_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"sigs.k8s.io/testing_frameworks/cluster"
	"sigs.k8s.io/testing_frameworks/dind"
	"sigs.k8s.io/testing_frameworks/lightweight"
)

var _ = Describe("Dind", func() {
	It("setup and teardown a dind cluster", func() {
		fixture := &dind.Dind{}
		fixture.Out = GinkgoWriter
		fixture.Err = GinkgoWriter

		config := cluster.Config{}
		config.API.BindURL = &url.URL{Scheme: "http", Host: "localhost:1234"}
		config.Shape.NodeCount = 3

		Expect(fixture.Setup(config)).To(Succeed())

		url := fixture.ClientConfig()
		Expect(url.Port()).To(Equal("1234"))

		kubectl := &lightweight.KubeCtl{}
		kubectl.Opts = append(kubectl.Opts, fmt.Sprintf("--server=%s", url.Host))

		stdout, stderr, err := kubectl.Run("get", "nodes", "-o", "json")

		// io.Copy(GinkgoWriter, stdout)
		io.Copy(GinkgoWriter, stderr)
		Expect(err).NotTo(HaveOccurred())

		nodes, err := parseNodes(stdout)
		Expect(err).NotTo(HaveOccurred())

		Expect(nodes.Items).To(HaveLen(4))
		_, workerCount := countNodes(nodes)
		Expect(workerCount).To(Equal(3))

		Expect(fixture.Teardown()).To(Succeed())
	})
})

func parseNodes(stdout io.Reader) (kubeNodes, error) {
	nodes := kubeNodes{}
	err := json.NewDecoder(stdout).Decode(&nodes)
	if err != nil {
		return kubeNodes{}, err
	}
	return nodes, nil
}

func countNodes(nodes kubeNodes) (int, int) {
	var workerCount int
	var masterCount int
	for _, node := range nodes.Items {
		if _, ok := node.Metadata.Labels["node-role.kubernetes.io/master"]; !ok {
			workerCount++
		} else {
			masterCount++
		}
	}
	return masterCount, workerCount
}

type kubeNodes struct {
	Items []struct {
		Metadata struct {
			Labels map[string]string
		}
	}
}
