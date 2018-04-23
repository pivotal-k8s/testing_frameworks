package integration

import (
	"net/url"

	"sigs.k8s.io/testing_frameworks/cluster"
	"sigs.k8s.io/testing_frameworks/integration/internal"
)

// Etcd knows how to run an etcd server.
type Etcd struct {
	// URL is the address the Etcd should listen on for client connections.
	//
	// If this is not specified, we default to a random free port on localhost.
	URL *url.URL

	// ClusterConfig is the kubeadm-compatible configuration for
	// clusters, which is partially supported by this framework.
	//
	// The elements of the ClusterConfig which are supported by
	// this framework are:
	//
	// - ClusterConfig.Etcd.DataDir
	// - ClusterConfig.Etcd.ExtraArgs
	ClusterConfig cluster.Config

	processState *internal.ProcessState
}

// Start starts the etcd, waits for it to come up, and returns an error, if one
// occoured.
func (e *Etcd) Start() error {
	var err error

	e.processState = &internal.ProcessState{}

	e.processState.DefaultedProcessInput, err = internal.DoDefaulting(
		"etcd",
		e.URL,
		e.ClusterConfig.Etcd.DataDir,
		e.ClusterConfig.Etcd.ProcessConfig.Path,
		e.ClusterConfig.Etcd.ProcessConfig.StartTimeout,
		e.ClusterConfig.Etcd.ProcessConfig.StopTimeout,
	)
	if err != nil {
		return err
	}

	e.processState.StartMessage = internal.GetEtcdStartMessage(e.processState.URL)

	e.URL = &e.processState.URL

	tmplData := struct {
		URL     *url.URL
		DataDir string
	}{
		e.URL,
		e.processState.Dir,
	}

	args := flattenArgs(e.ClusterConfig.Etcd.ExtraArgs)

	e.processState.Args, err = internal.RenderTemplates(
		internal.DoEtcdArgDefaulting(args), tmplData,
	)
	if err != nil {
		return err
	}

	return e.processState.Start(
		e.ClusterConfig.Etcd.ProcessConfig.Out,
		e.ClusterConfig.Etcd.ProcessConfig.Err,
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
