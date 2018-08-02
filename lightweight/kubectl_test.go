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

package lightweight_test

import (
	"io/ioutil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "sigs.k8s.io/testing_frameworks/lightweight"
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
})
