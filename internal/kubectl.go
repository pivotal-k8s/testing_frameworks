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
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"

	"sigs.k8s.io/testing_frameworks/cluster"
	"sigs.k8s.io/testing_frameworks/cluster/type/base"
)

// KubeCtl is a wrapper around the kubectl binary.
type KubeCtl struct {
	// Path where the kubectl binary can be found.
	//
	// If this is left empty, we will attempt to locate a binary, by checking for
	// the TEST_ASSET_KUBECTL environment variable, and the default test assets
	// directory. See the "Binaries" section above (in doc.go) for details.
	Path string

	// Opts can be used to configure additional flags which will be used each
	// time the wrapped binary is called.
	//
	// For example, you might want to use this to set the URL of the APIServer to
	// connect to.
	Opts []string

	// KubeConfig holds the information of how to connect to a cluster. The
	// format is compatible with `clientcmd.Config` when serialized.
	KubeConfig *base.Config
}

// Run executes the wrapped binary with some preconfigured options and the
// arguments given to this method. It returns Readers for the stdout and
// stderr.
func (k *KubeCtl) Run(args ...string) (stdout, stderr io.Reader, err error) {
	if k.Path == "" {
		k.Path = BinPathFinder("lightweight", "kubectl")
	}

	stdoutBuffer := &bytes.Buffer{}
	stderrBuffer := &bytes.Buffer{}

	allArgs := append(k.Opts, args...)

	cmd := exec.Command(k.Path, allArgs...)
	cmd.Stdout = stdoutBuffer
	cmd.Stderr = stderrBuffer

	if k.KubeConfig != nil {
		cleanup, err := k.configureKubeConfig(cmd)
		if err != nil {
			return nil, nil, err
		}

		defer cleanup()
	}

	err = cmd.Run()

	return stdoutBuffer, stderrBuffer, err
}

func (k *KubeCtl) Configure(f cluster.Fixture) *KubeCtl {
	c := f.ClientConfig()
	k.KubeConfig = &c
	return k
}

func (k *KubeCtl) configureKubeConfig(cmd *exec.Cmd) (func(), error) {
	confFile, err := ioutil.TempFile("", ".kubeconf-")
	if err != nil {
		return nil, err
	}

	if err := json.NewEncoder(confFile).Encode(k.KubeConfig); err != nil {
		return nil, err
	}

	cmd.Env = append(cmd.Env, fmt.Sprintf("KUBECONFIG=%s", confFile.Name()))

	return func() {
		os.Remove(confFile.Name())
	}, nil
}
