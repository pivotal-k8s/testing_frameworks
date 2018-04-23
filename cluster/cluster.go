// Package cluster is a Test Cluster Framework -- see the motivating
// document here:
// https://docs.google.com/document/d/13bMjmWpsdkgbY-JayrcU-e_QNwRJCP-rHjtqdeeoQHo/edit?ts=5aa005c9#heading=h.75awtuvlo3ad
//
// This package aims to provide a consistent abstraction over a number
// of different ways of creating kubernetes clusters for testing.
//
// To make your test cluster compatible with this framework, just
// implement the Fixture interface.  If you require more config
// options than we have so far, add them to the Config struct. Bear in
// mind while doing so that we aim to maintain compatibility with
// kubeadm, so if your new config options describe properties that
// already exist in the kubeadm config, you should use
// kubeadm-compatibile terminology.
package cluster

import (
	"io"
	"net/url"
	"time"
)

// Fixture is some kind of test cluster fixture, which can be started, interacted with, and stopped.
type Fixture interface {
	// Setup starts the test cluster according to the provided
	// configuration. If the config asks for a feature that is not
	// supported by this Fixture implementation, or if the test
	// cluster fails to start, return an error.
	//
	// This should block until the test cluster control plane is
	// ready to recieve client connections.
	Setup(config Config) error

	// TearDown cleanly stops the test cluster. If we can't stop
	// cleanly, return an error.
	//
	// TearDown should block until the test cluster has stopped,
	// or we have given up on stopping it and returned an error.
	//
	// TearDown should be idempotent. If a user calls TearDown
	// twice in a row, and the first call succeded, then the
	// second call should also succeed.
	TearDown() error

	// ClientConfig returns the URL at which you can find the APIServer
	// of the test cluster. In future, this will likely become a
	// rest.Config from
	// https://github.com/kubernetes/client-go/blob/master/rest/config.go.
	ClientConfig() *url.URL
}

// KubeadmMasterConfig is a struct which is used as a nested struct in `Config`.
//
// We use this indirection to be prepared if we'd be vendored into k/k -- then
// we could remove `KubeadmMasterConfig` and use the actual one from within k/k.
type KubeadmMasterConfig struct {
	// CertificatesDir specifies where to store or look for all required certificates.
	// See also https://github.com/kubernetes/kubernetes/blob/c8cded58d71e36665bd345a70fbe404e7523abb8/cmd/kubeadm/app/apis/kubeadm/types.go#L104
	CertificatesDir string

	// APIServerExtraArgs is a set of extra flags to pass to the API Server or override
	// default ones in form of <flagname>=<value>.
	// See also https://github.com/kubernetes/kubernetes/blob/c8cded58d71e36665bd345a70fbe404e7523abb8/cmd/kubeadm/app/apis/kubeadm/types.go#L81
	APIServerExtraArgs map[string]string
}

// LightWeightMasterConfig is a struct used as a nested struct in `Config` to
// add aditional configuration properties needed by the "lightweight"
// implementation to the main `Config` struct.
type LightWeightMasterConfig struct {
	// APIServerProcessConfig hold configuration properties related to the
	// APIServer process.
	APIServerProcessConfig ProcessConfig
}

// Config is a struct into which you can parse a YAML or JSON config
// file (which should always be compatible with
// https://kubernetes.io/docs/reference/setup-tools/kubeadm/kubeadm-init/#config-file
// ) to describe your test cluster.
//
// To maintain compatibility with kubeadm, we should follow the
// patterns established in
// https://github.com/kubernetes/kubernetes/blob/c8cded58d71e36665bd345a70fbe404e7523abb8/cmd/kubeadm/app/apis/kubeadm/types.go#L30
type Config struct {
	// Etcd holds configuration for etcd.
	Etcd Etcd

	// API holds configuration for the k8s apiserver.
	API API

	// KubeadmMasterConfig is the nested struct holding all the configuration
	// supported by kubeadm
	KubeadmMasterConfig

	// LightWeightMasterConfig is the nested struct holding all the configuration
	// additionally supported by the "lightweight" framework
	LightWeightMasterConfig
}

// KubeadmEtcd is a struct which is used as a nested struct in `Etcd`.
//
// We use this indirection to be prepared if we'd be vendored into k/k -- then
// we could remove `KubeadmEtcd` and use the actual one from within k/k.
type KubeadmEtcd struct {
	// DataDir is the directory etcd will place its data.
	// Defaults to "/var/lib/etcd".
	DataDir string

	// ExtraArgs are extra arguments provided to the etcd binary
	// when run inside a static pod.
	ExtraArgs map[string]string
}

// LightWeightEtcd is a struct used as a nested struct in `Etcd` to
// add aditional configuration properties needed by the "lightweight"
// implementation to the main `Etcd` struct.
type LightWeightEtcd struct {
	// ProcessConfig holds configuration properties releated to the Etcd progress
	ProcessConfig ProcessConfig

	// BindURL is the URL Etcd should bind to
	BindURL *url.URL
}

// Etcd contains elements describing Etcd configuration.
// See also https://github.com/kubernetes/kubernetes/blob/c8cded58d71e36665bd345a70fbe404e7523abb8/cmd/kubeadm/app/apis/kubeadm/types.go#L163
type Etcd struct {
	// KubeadmEtcd is the nexted struct holding all the configuration for etcd
	// supported by kubeadm
	KubeadmEtcd

	// LightWeightEtcd is the nested struct holding all the configuration
	// additionally supported by the "lightweight" framework
	LightWeightEtcd
}

// KubeadmAPI is a struct which is used as a nested struct in `API`.
//
// We use this indirection to be prepared if we'd be vendored into k/k -- then
// we could remove `KubeadmEtcd` and use the actual one from within k/k.
//
// For now this is empty, as none of the currently supported framework supports
// any of these settings.
type KubeadmAPI struct{}

// LightWeightAPI is a struct used as a nested struct in `API` to
// add aditional configuration properties needed by the "lightweight"
// implementation to the main `API` struct.
type LightWeightAPI struct {
	// BindURL is a URL the API should listen on.
	//
	// If this is kept empty, it will be defaulted to a free port on "localhost".
	BindURL *url.URL
}

// API struct contains elements of API server address.
type API struct {
	// KubeadmAPI is the nested struct holding all the configuration for the API
	// server address supported by kubeadm.
	KubeadmAPI

	// LightWeightAPI is the nested struct holding all addirional configuration
	// for the API server adress supported by the "lightweight" framework.
	LightWeightAPI
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
