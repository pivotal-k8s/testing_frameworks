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

package lightweight

import (
	"fmt"
	"net/url"

	"sigs.k8s.io/testing_frameworks/cluster"
	"sigs.k8s.io/testing_frameworks/cluster/type/base"
)

// ControlPlane is a struct that knows how to start your test control plane.
//
// Right now, that means Etcd and your APIServer. This is likely to increase in
// future.
type ControlPlane struct {
	APIServer *APIServer
	Etcd      *Etcd
	Nodes     []*Node
}

// Setup will start your control plane processes according to the
// supplied configuration.
func (f *ControlPlane) Setup(config cluster.Config) error {
	f.Etcd = &Etcd{}
	f.Etcd.ClusterConfig = config

	if f.APIServer == nil {
		f.APIServer = &APIServer{}
	}
	f.APIServer.ClusterConfig = config

	for i := 1; i <= config.Shape.NodeCount; i++ {
		f.Nodes = append(f.Nodes, &Node{
			ClusterConfig: config,
		})
	}

	return f.Start()
}

// Start will start your control plane processes. To stop them, call Stop().
func (f *ControlPlane) Start() error {
	if f.Etcd == nil {
		f.Etcd = &Etcd{}
	}
	if err := f.Etcd.Start(); err != nil {
		return err
	}

	if f.APIServer == nil {
		f.APIServer = &APIServer{}
	}
	f.APIServer.EtcdURL = &f.Etcd.processState.URL
	if err := f.APIServer.Start(); err != nil {
		return err
	}

	for _, n := range f.Nodes {
		if err := n.Start(); err != nil {
			return err
		}
	}

	return nil
}

// TearDown will stop your control plane processes. This is an alias for Stop()
func (f *ControlPlane) TearDown() error {
	return f.Stop()
}

// Stop will stop your control plane processes, and clean up their data.
func (f *ControlPlane) Stop() error {
	for _, node := range f.Nodes {
		if err := node.Stop(); err != nil {
			return err
		}
	}
	if f.APIServer != nil {
		if err := f.APIServer.Stop(); err != nil {
			return err
		}
	}
	if f.Etcd != nil {
		if err := f.Etcd.Stop(); err != nil {
			return err
		}
	}
	return nil
}

// ClientConfig returns all the configuration a client needs to connect to this
// cluster's kubernetes API.
func (f *ControlPlane) ClientConfig() base.Config {
	c := base.Config{}

	cluster := base.NamedCluster{}
	cluster.Name = "lightweight-cluster"
	cluster.Cluster.Server = f.APIURL().String()

	ctx := base.NamedContext{}
	ctx.Name = "lightweight-context"
	ctx.Context.Cluster = cluster.Name

	c.Clusters = []base.NamedCluster{cluster}
	c.Contexts = []base.NamedContext{ctx}
	c.CurrentContext = ctx.Name

	return c
}

// APIURL returns the URL you should connect to to talk to your API.
func (f *ControlPlane) APIURL() *url.URL {
	return &f.APIServer.processState.URL
}

// KubeCtl returns a pre-configured KubeCtl, ready to connect to this
// ControlPlane.
func (f *ControlPlane) KubeCtl() *KubeCtl {
	k := &KubeCtl{}
	k.Opts = append(k.Opts, fmt.Sprintf("--server=%s", f.APIURL()))
	return k
}
