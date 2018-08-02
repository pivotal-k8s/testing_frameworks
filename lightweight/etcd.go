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
	"net/url"

	"sigs.k8s.io/testing_frameworks/cluster"
	"sigs.k8s.io/testing_frameworks/lightweight/internal"
)

// Etcd knows how to run an etcd server.
type Etcd struct {
	// ClusterConfig is the kubeadm-compatible configuration for
	// clusters, which is partially supported by this framework.
	ClusterConfig cluster.Config

	processState *internal.ProcessState
}

// Start starts the etcd, waits for it to come up, and returns an error, if one
// occoured.
func (e *Etcd) Start() error {
	var err error

	clusterConf := cluster.DoDefaulting(e.ClusterConfig)

	e.processState = &internal.ProcessState{}

	e.processState.DefaultedProcessInput, err = internal.DoDefaulting(
		"etcd",
		clusterConf.Etcd.BindURL,
		clusterConf.Etcd.Local.DataDir,
		clusterConf.Etcd.ProcessConfig.Path,
		clusterConf.Etcd.ProcessConfig.StartTimeout,
		clusterConf.Etcd.ProcessConfig.StopTimeout,
	)
	if err != nil {
		return err
	}

	e.processState.StartMessage = internal.GetEtcdStartMessage(e.processState.URL)

	tmplData := struct {
		URL     *url.URL
		DataDir string
	}{
		&e.processState.URL,
		e.processState.Dir,
	}

	args := flattenArgs(clusterConf.Etcd.Local.ExtraArgs)

	e.processState.Args, err = internal.RenderTemplates(
		internal.DoEtcdArgDefaulting(args), tmplData,
	)
	if err != nil {
		return err
	}

	return e.processState.Start(
		clusterConf.Etcd.ProcessConfig.Out,
		clusterConf.Etcd.ProcessConfig.Err,
	)
}

// Stop stops this process gracefully, waits for its termination, and cleans up
// the DataDir if necessary.
func (e *Etcd) Stop() error {
	return e.processState.Stop()
}

// EtcdDefaultArgs exposes the default args for Etcd so that you
// can use those to append your own additional arguments.
//
// The internal default arguments are explicitely copied here, we don't want to
// allow users to change the internal ones.
var EtcdDefaultArgs = append([]string{}, internal.EtcdDefaultArgs...)
