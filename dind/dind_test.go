/*
Copyright 2018 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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
	var (
		fixture *dind.Dind
	)

	AfterEach(func() {
		Expect(fixture.TearDown()).To(Succeed())
	})

	It("sets up and tears down a dind cluster with default settings", func() {
		fixture = &dind.Dind{}
		fixture.Out = GinkgoWriter
		fixture.Err = GinkgoWriter

		config := cluster.Config{}
		Expect(fixture.Setup(config)).To(Succeed())
	})

	It("setup and teardown a dind cluster", func() {
		fixture = &dind.Dind{}
		fixture.Out = GinkgoWriter
		fixture.Err = GinkgoWriter

		config := cluster.Config{}
		config.API.BindURL = &url.URL{Scheme: "http", Host: "localhost:1234"}
		config.Shape.NodeCount = 2

		Expect(fixture.Setup(config)).To(Succeed())

		url := fixture.ClientConfig()
		Expect(url.Port()).To(Equal("1234"))

		kubectl := &dindKubeCtl{URL: fixture.ClientConfig()}
		stdout, _, err := kubectl.Run("get", "nodes", "-o", "json")
		Expect(err).NotTo(HaveOccurred())

		nodes, err := parseNodes(stdout)
		Expect(err).NotTo(HaveOccurred())
		Expect(nodes.Items).To(HaveLen(3))

		_, workerCount := countNodes(nodes)
		Expect(workerCount).To(Equal(2))
	})
})

// TODO: move lightweight's kubectl into a shared package
type dindKubeCtl struct {
	URL     *url.URL
	kubectl *lightweight.KubeCtl
}

func (k *dindKubeCtl) Run(args ...string) (io.Reader, io.Reader, error) {
	if k.kubectl == nil {
		k.kubectl = &lightweight.KubeCtl{}
	}

	k.kubectl.Opts = append(
		k.kubectl.Opts,
		fmt.Sprintf("--server=%s", k.URL),
	)

	stdout, stderr, err := k.kubectl.Run(args...)
	stdout = io.TeeReader(stdout, GinkgoWriter)
	stderr = io.TeeReader(stderr, GinkgoWriter)

	return stdout, stderr, err
}

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
