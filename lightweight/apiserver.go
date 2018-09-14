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
	"sigs.k8s.io/testing_frameworks/lightweight/internal"
)

// APIServer knows how to run a kubernetes apiserver.
type APIServer struct {
	// ClusterConfig is the kubeadm-compatible configuration for
	// clusters, which is partially supported by this framework.
	ClusterConfig cluster.Config

	// EtcdURL is the URL of the Etcd the APIServer should use.
	//
	// If this is not specified, the Start() method will return an error.
	EtcdURL *url.URL

	processState *internal.ProcessState
}

// Start starts the apiserver, waits for it to come up, and returns an error,
// if occurred.
func (s *APIServer) Start() error {
	if s.EtcdURL == nil {
		return fmt.Errorf("expected EtcdURL to be configured")
	}

	var err error

	clusterConf := cluster.DoDefaulting(s.ClusterConfig)

	s.processState = &internal.ProcessState{}

	apiURL, err := url.Parse(clusterConf.ControlPlaneEndpoint)
	if err != nil {
		return err
	}
	s.processState.DefaultedProcessInput, err = internal.DoDefaulting(
		"kube-apiserver",
		apiURL,
		clusterConf.CertificatesDir,
		clusterConf.APIServerProcessConfig.Path,
		clusterConf.APIServerProcessConfig.StartTimeout,
		clusterConf.APIServerProcessConfig.StopTimeout,
	)
	if err != nil {
		return err
	}

	s.processState.HealthCheckEndpoint = "/healthz"

	tmplData := struct {
		EtcdURL *url.URL
		URL     *url.URL
		CertDir string
	}{
		s.EtcdURL,
		&s.processState.URL,
		s.processState.Dir,
	}

	args := flattenArgs(clusterConf.APIServerExtraArgs)

	s.processState.Args, err = internal.RenderTemplates(
		internal.DoAPIServerArgDefaulting(args), tmplData,
	)
	if err != nil {
		return err
	}

	return s.processState.Start(
		clusterConf.APIServerProcessConfig.Out,
		clusterConf.APIServerProcessConfig.Err,
	)
}

// Stop stops this process gracefully, waits for its termination, and cleans up
// the CertDir if necessary.
func (s *APIServer) Stop() error {
	return s.processState.Stop()
}

// APIServerDefaultArgs exposes the default args for the APIServer so that you
// can use those to append your own additional arguments.
//
// The internal default arguments are explicitely copied here, we don't want to
// allow users to change the internal ones.
var APIServerDefaultArgs = append([]string{}, internal.APIServerDefaultArgs...)

func flattenArgs(mappedArgs map[string]string) []string {
	args := []string{}
	for k, v := range mappedArgs {
		if v == "" {
			args = append(args, k)
		} else {
			args = append(args, fmt.Sprintf("%s=%s", k, v))
		}
	}
	return args
}
