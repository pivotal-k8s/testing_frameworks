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

package integration_tests

import (
	"io/ioutil"
	"path/filepath"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"sigs.k8s.io/testing_frameworks/cluster"
	"sigs.k8s.io/testing_frameworks/cluster/type/base"
	"sigs.k8s.io/testing_frameworks/lightweight"
)

var _ = Describe("Cluster Framework Compliance", func() {
	var fixture cluster.Fixture

	AfterEach(func() {
		Expect(fixture.TearDown()).To(Succeed())
	})

	It("Successfully manages the control plane lifecycle", func() {
		var err error

		fixture = &lightweight.ControlPlane{}
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

		fixture = &lightweight.ControlPlane{}

		config := cluster.Config{}
		config.Etcd.Local = &base.LocalEtcd{
			DataDir: dataDir,
		}

		err = fixture.Setup(config)
		Expect(err).NotTo(HaveOccurred())

		Expect(dataDir).To(BeADirectory())

		Expect(fixture.TearDown()).To(Succeed())
	})

	It("Manages a configured apiserver certificate directory", func() {
		dir, err := ioutil.TempDir("", "")
		Expect(err).NotTo(HaveOccurred())

		certDir := filepath.Join(dir, "apiserver-cert-dir")
		Expect(certDir).NotTo(BeAnExistingFile())

		fixture = &lightweight.ControlPlane{}

		config := cluster.Config{}
		config.CertificatesDir = certDir

		err = fixture.Setup(config)
		Expect(err).NotTo(HaveOccurred())

		Expect(certDir).To(BeADirectory())

		Expect(fixture.TearDown()).To(Succeed())
	})

	It("Fails on an unknown commandline argument", func() {
		stdErr := gbytes.NewBuffer()

		config := cluster.Config{}
		config.APIServerProcessConfig.Err = stdErr
		config.APIServerProcessConfig.StartTimeout = 500 * time.Millisecond
		config.APIServerExtraArgs = map[string]string{
			"--some-silly-arg": "",
		}

		fixture = &lightweight.ControlPlane{
			APIServer: &lightweight.APIServer{
				ClusterConfig: config,
			},
		}

		Expect(fixture.Setup(config)).NotTo(Succeed())
		Expect(stdErr).To(gbytes.Say("some-silly-arg"))
	})

	It("Supports a shape with multiple node sets", func() {
		config := cluster.Config{}
		config.Shape.NodeCount = 2

		fixture = &lightweight.ControlPlane{}

		Expect(fixture.Setup(config)).To(Succeed())
	})
})
