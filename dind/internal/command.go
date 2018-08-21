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

package internal

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"os/exec"

	"sigs.k8s.io/testing_frameworks/cluster"
	"sigs.k8s.io/testing_frameworks/internal"
)

func UpCommand(label string, stdOut, stdErr io.Writer, clusterConfig cluster.Config) *exec.Cmd {
	additionalEnv := []string{
		fmt.Sprintf("NUM_NODES=%d", clusterConfig.Shape.NodeCount),
	}

	if ok, port := getPortFromURL(clusterConfig.API.BindURL); ok {
		additionalEnv = append(additionalEnv, fmt.Sprintf("APISERVER_PORT=%s", port))
	}

	if v := clusterConfig.KubernetesVersion; v != "" {
		additionalEnv = append(
			additionalEnv,
			fmt.Sprintf("DIND_IMAGE=mirantis/kubeadm-dind-cluster:%s", v),
		)
	}

	cmd := clusterCmd(label, "up", additionalEnv...)
	return attachIO(cmd, stdOut, stdErr)
}

func CleanCommand(label string, stdOut, stdErr io.Writer) *exec.Cmd {
	cmd := clusterCmd(label, "clean")
	return attachIO(cmd, stdOut, stdErr)
}

func APIServerPortCommand(label string, stdOut, stdErr io.Writer) *exec.Cmd {
	cmd := clusterCmd(label, "apiserver-port")
	return attachIO(cmd, stdOut, stdErr)
}

func clusterCmd(clusterLabel string, subCommand string, additionalEnvs ...string) *exec.Cmd {
	binPath := internal.BinPathFinder("dind", "dind-cluster.sh")

	cmd := exec.Command(binPath, subCommand) // #nosec
	cmd.Env = clusterEnv(clusterLabel, additionalEnvs...)

	return cmd
}

func attachIO(cmd *exec.Cmd, stdOut, stdErr io.Writer) *exec.Cmd {
	if stdOut != nil {
		cmd.Stdout = stdOut
	}
	if stdErr != nil {
		cmd.Stderr = stdErr
	}
	return cmd
}

func clusterEnv(label string, additionalEnv ...string) []string {
	env := append(os.Environ(),
		fmt.Sprintf("DIND_LABEL=%s", label),
	)
	for _, e := range additionalEnv {
		env = append(env, e)
	}
	return env
}

func getPortFromURL(u *url.URL) (bool, string) {
	if u == nil {
		return false, ""
	}
	return true, u.Port()
}