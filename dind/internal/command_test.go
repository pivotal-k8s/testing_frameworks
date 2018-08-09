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
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"

	"sigs.k8s.io/testing_frameworks/cluster"
	. "sigs.k8s.io/testing_frameworks/dind/internal"
)

var _ = Describe("command", func() {
	var cmd *exec.Cmd
	var clusterConfig cluster.Config

	BeforeEach(func() {
		clusterConfig = cluster.Config{}
		clusterConfig.Shape.NodeCount = 1234
	})

	Context("UpCommand", func() {
		JustBeforeEach(func() {
			cmd = UpCommand("some_label", nil, nil, clusterConfig)
		})

		It("has minimal setup", func() {
			Expect(cmd).To(haveMinimalSetup())
		})
		It("has subCommand set", func() {
			Expect(cmd.Args[1]).To(Equal("up"))
		})
		It("sets NUM_NODES", func() {
			Expect(cmd.Env).To(haveVariableWithValue("NUM_NODES", "1234"))
		})

		It("doesn't set APISERVER_PORT by default", func() {
			Expect(cmd.Env).NotTo(haveVariable("APISERVER_PORT"))
		})
		Context("with APIServer URL configured", func() {
			BeforeEach(func() {
				clusterConfig.API.BindURL = &url.URL{Host: ":5678"}
			})
			It("sets APISERVER_PORT", func() {
				Expect(cmd.Env).To(haveVariableWithValue("APISERVER_PORT", "5678"))
			})
		})

		It("doesn't set DIND_IMAGE by default", func() {
			Expect(cmd.Env).NotTo(haveVariable("DIND_IMAGE"))
		})
		Context("with kubernetes version configured", func() {
			BeforeEach(func() {
				clusterConfig.KubernetesVersion = "some_version"
			})
			It("sets DIND_IMAGE accordingly", func() {
				Expect(cmd.Env).To(haveVariableWithValue("DIND_IMAGE", "mirantis/kubeadm-dind-cluster:some_version"))
			})
		})
	})

	Context("CleanCommand", func() {
		JustBeforeEach(func() {
			cmd = CleanCommand("some_label", nil, nil)
		})

		It("has minimal setup", func() {
			Expect(cmd).To(haveMinimalSetup())
		})
		It("has subCommand set", func() {
			Expect(cmd.Args[1]).To(Equal("clean"))
		})
	})

	Context("APIServerPortCommand", func() {
		JustBeforeEach(func() {
			cmd = APIServerPortCommand("some_label", nil, nil)
		})

		It("has minimal setup", func() {
			Expect(cmd).To(haveMinimalSetup())
		})
		It("has subCommand set", func() {
			Expect(cmd.Args[1]).To(Equal("apiserver-port"))
		})
	})
})

func haveMinimalSetup() types.GomegaMatcher {
	getPath := func(cmd *exec.Cmd) string {
		return cmd.Path
	}
	getEnv := func(cmd *exec.Cmd) []string {
		return cmd.Env
	}
	return And(
		Not(WithTransform(getPath, BeEmpty())),
		WithTransform(getEnv, haveNonEmptyVariable("DIND_LABEL")),
	)
}

func haveNonEmptyVariable(name string) types.GomegaMatcher {
	return ContainElement(
		MatchRegexp("^%s=[^\\s]+", name),
	)
}

func haveVariableWithValue(name, value string) types.GomegaMatcher {
	return ContainElement(
		MatchRegexp("^%s=%s$", name, value),
	)
}

func haveVariable(name string) types.GomegaMatcher {
	return ContainElement(HavePrefix("%s=", name))
}
