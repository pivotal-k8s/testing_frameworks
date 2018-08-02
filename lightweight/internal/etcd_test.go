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
	"net/url"

	. "sigs.k8s.io/testing_frameworks/lightweight/internal"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Etcd", func() {
	It("defaults Args if they are empty", func() {
		initialArgs := []string{}
		defaultedArgs := DoEtcdArgDefaulting(initialArgs)
		Expect(defaultedArgs).To(BeEquivalentTo(EtcdDefaultArgs))
	})

	It("keeps Args as is if they are not empty", func() {
		initialArgs := []string{"--eins", "--zwei=2"}
		defaultedArgs := DoEtcdArgDefaulting(initialArgs)
		Expect(defaultedArgs).To(BeEquivalentTo([]string{
			"--eins", "--zwei=2",
		}))
	})
})

var _ = Describe("GetEtcdStartMessage()", func() {
	Context("when using a non tls URL", func() {
		It("generates valid start message", func() {
			url := url.URL{
				Scheme: "http",
				Host:   "some.insecure.host:1234",
			}
			message := GetEtcdStartMessage(url)
			Expect(message).To(Equal("serving insecure client requests on some.insecure.host"))
		})
	})
	Context("when using a tls URL", func() {
		It("generates valid start message", func() {
			url := url.URL{
				Scheme: "https",
				Host:   "some.secure.host:8443",
			}
			message := GetEtcdStartMessage(url)
			Expect(message).To(Equal("serving client requests on some.secure.host"))
		})
	})
})
