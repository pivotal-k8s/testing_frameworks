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
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "sigs.k8s.io/testing_frameworks/lightweight/internal"
)

var _ = Describe("Apiserver", func() {
	It("defaults Args if they are empty", func() {
		initialArgs := []string{}
		defaultedArgs := DoAPIServerArgDefaulting(initialArgs)
		Expect(defaultedArgs).To(BeEquivalentTo(APIServerDefaultArgs))
	})

	It("keeps Args as is if they are not empty", func() {
		initialArgs := []string{"--one", "--two=2"}
		defaultedArgs := DoAPIServerArgDefaulting(initialArgs)
		Expect(defaultedArgs).To(BeEquivalentTo([]string{
			"--one", "--two=2",
		}))
	})
})
