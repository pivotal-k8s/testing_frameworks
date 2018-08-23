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

package internal_test

import (
	"io"
	"io/ioutil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"sigs.k8s.io/testing_frameworks/cluster"
	"sigs.k8s.io/testing_frameworks/cluster/type/base"
	. "sigs.k8s.io/testing_frameworks/internal"
)

var _ = Describe("Kubectl", func() {
	It("runs kubectl", func() {
		k := &KubeCtl{Path: "bash"}
		args := []string{"-c", "echo 'something'"}
		stdout, stderr, err := k.Run(args...)
		Expect(err).NotTo(HaveOccurred())
		Expect(stdout).To(ContainSubstring("something"))
		bytes, err := ioutil.ReadAll(stderr)
		Expect(err).NotTo(HaveOccurred())
		Expect(bytes).To(BeEmpty())
	})

	Context("when the command returns a non-zero exit code", func() {
		It("returns an error", func() {
			k := &KubeCtl{Path: "bash"}
			args := []string{
				"-c", "echo 'this is StdErr' >&2; echo 'but this is StdOut' >&1; exit 66",
			}

			stdout, stderr, err := k.Run(args...)

			Expect(err).To(MatchError(ContainSubstring("exit status 66")))

			Expect(stdout).To(ContainSubstring("but this is StdOut"))
			Expect(stderr).To(ContainSubstring("this is StdErr"))
		})
	})

	Context("when configured for a fixture", func() {
		It("holds some configuration", func() {
			f := testFixture{}
			k := &KubeCtl{Path: "bash"}
			a := []string{"-c", "echo $KUBECONFIG"}

			Expect(k.KubeConfig).To(BeNil())

			k.Configure(f)

			Expect(k.KubeConfig).NotTo(BeNil())

			o, e, err := k.Run(a...)
			Expect(err).NotTo(HaveOccurred())

			Expect(asString(o)).NotTo(BeEmpty())
			Expect(asString(e)).To(BeEmpty())
		})
	})
})

type testFixture struct{}

func (tf testFixture) Setup(c cluster.Config) error {
	return nil
}

func (tf testFixture) TearDown() error {
	return nil
}

func (tf testFixture) ClientConfig() base.Config {
	return base.Config{}
}

func asString(r io.Reader) string {
	b, err := ioutil.ReadAll(r)
	ExpectWithOffset(1, err).NotTo(HaveOccurred())
	return string(b)
}
