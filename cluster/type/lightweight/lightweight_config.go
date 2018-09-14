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

// Package lightweight holds the configuration extensions for the lightweight
// test cluster implementation.
package lightweight

import (
	"io"
	"net/url"
	"time"
)

// MasterConfiguration is a struct used as a nested struct in `Config` to
// add additional configuration properties needed by the "lightweight"
// implementation to the main `Config` struct.
type MasterConfiguration struct {
	// APIServerProcessConfig hold configuration properties related to the
	// APIServer process.
	APIServerProcessConfig ProcessConfig
}

// Etcd is a struct used as a nested struct in `Etcd` to
// add additional configuration properties needed by the "lightweight"
// implementation to the main `Etcd` struct.
type Etcd struct {
	// ProcessConfig holds configuration properties releated to the Etcd progress
	ProcessConfig ProcessConfig

	// BindURL is the URL Etcd should bind to.
	BindURL *url.URL
}

// APIEndpoint is a struct used as a nested struct in `API` to
// add additional configuration properties needed by the "lightweight"
// implementation to the main `API` struct.
type APIEndpoint struct {
	// BindURL is a URL the API should listen on.
	//
	// If this is kept empty, it will be defaulted to a free port on "localhost".
	BindURL *url.URL
}

// ProcessConfig is configuring certain properties for processes.
type ProcessConfig struct {
	// Path is the path to the binary.
	//
	// If this is left as the empty string, we will attempt to locate a binary,
	// by checking for the TEST_ASSET_KUBE_<COMPONENT> (e.g.
	// TEST_ASSET_KUBE_APISERVER, TEST_ASSET_KUBE_ETCD, ...) environment
	// variable, and the default test assets directory. See the "Binaries"
	// section above (in doc.go) for details.
	Path string

	// StartTimeout, StopTimeout specify the time the the process is allowed to
	// take when starting and stoppping before an error is emitted.
	//
	// If not specified, these default to 20 seconds.
	StartTimeout time.Duration
	StopTimeout  time.Duration

	// Out, Err specify where the process should write its StdOut, StdErr to.
	//
	// If not specified, the output will be discarded.
	Out io.Writer
	Err io.Writer
}
