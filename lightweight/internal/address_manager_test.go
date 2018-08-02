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
	. "sigs.k8s.io/testing_frameworks/lightweight/internal"

	"fmt"
	"net"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("AddressManager", func() {
	var addressManager *AddressManager
	BeforeEach(func() {
		addressManager = &AddressManager{}
	})

	Describe("Initialize", func() {
		It("returns a free port and an address to bind to", func() {
			port, host, err := addressManager.Initialize()

			Expect(err).NotTo(HaveOccurred())
			Expect(host).To(Equal("127.0.0.1"))
			Expect(port).NotTo(Equal(0))

			addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", host, port))
			Expect(err).NotTo(HaveOccurred())
			l, err := net.ListenTCP("tcp", addr)
			defer func() {
				Expect(l.Close()).To(Succeed())
			}()
			Expect(err).NotTo(HaveOccurred())
		})

		Context("initialized multiple times", func() {
			It("fails", func() {
				_, _, err := addressManager.Initialize()
				Expect(err).NotTo(HaveOccurred())
				_, _, err = addressManager.Initialize()
				Expect(err).To(MatchError(ContainSubstring("already initialized")))
			})
		})
	})
	Describe("Port", func() {
		It("returns an error if Initialize has not been called yet", func() {
			_, err := addressManager.Port()
			Expect(err).To(MatchError(ContainSubstring("not initialized yet")))
		})
		It("returns the same port as previously allocated by Initialize", func() {
			expectedPort, _, err := addressManager.Initialize()
			Expect(err).NotTo(HaveOccurred())
			actualPort, err := addressManager.Port()
			Expect(err).NotTo(HaveOccurred())
			Expect(actualPort).To(Equal(expectedPort))
		})
	})
	Describe("Host", func() {
		It("returns an error if Initialize has not been called yet", func() {
			_, err := addressManager.Host()
			Expect(err).To(MatchError(ContainSubstring("not initialized yet")))
		})
		It("returns the same port as previously allocated by Initialize", func() {
			_, expectedHost, err := addressManager.Initialize()
			Expect(err).NotTo(HaveOccurred())
			actualHost, err := addressManager.Host()
			Expect(err).NotTo(HaveOccurred())
			Expect(actualHost).To(Equal(expectedHost))
		})
	})
})
