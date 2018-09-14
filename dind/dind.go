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

package dind

import (
	"bytes"
	"fmt"
	"io"

	"sigs.k8s.io/testing_frameworks/cluster"
	"sigs.k8s.io/testing_frameworks/cluster/type/base"
	"sigs.k8s.io/testing_frameworks/dind/internal"
)

type Dind struct {
	Out io.Writer
	Err io.Writer

	label string
}

func (d *Dind) Setup(c cluster.Config) error {
	label, err := internal.RandomString(10)
	if err != nil {
		return err
	}
	d.label = label

	cmd, err := internal.UpCommand(d.label, d.Out, d.Err, c)
	if err != nil {
		return err
	}

	return cmd.Run()
}

func (d *Dind) TearDown() error {
	return internal.
		CleanCommand(d.label, d.Out, d.Err).
		Run()
}

func (d *Dind) ClientConfig() base.Config {
	c := base.Config{}

	// TODO: let that error bubble up
	APIServerPort, _ := d.getAPIServerPort()

	cluster := base.NamedCluster{}
	cluster.Name = "k-d-c"
	cluster.Cluster.Server = fmt.Sprintf("http://localhost:%s", APIServerPort)

	ctx := base.NamedContext{}
	ctx.Name = "k-d-c"
	ctx.Context.Cluster = cluster.Name

	c.Clusters = []base.NamedCluster{cluster}
	c.Contexts = []base.NamedContext{ctx}
	c.CurrentContext = ctx.Name

	return c
}

func (d *Dind) getAPIServerPort() (string, error) {
	stdout := &bytes.Buffer{}
	cmd := internal.APIServerPortCommand(d.label, stdout, nil)

	if err := cmd.Run(); err != nil {
		return "", err
	}

	var port int
	_, err := fmt.Fscanf(stdout, "%d\n", &port)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%d", port), nil
}
