package integration

import (
	"fmt"
	"io"
	"net/url"
	"time"

	"sigs.k8s.io/testing_frameworks/cluster"
	"sigs.k8s.io/testing_frameworks/integration/internal"
)

// APIServer knows how to run a kubernetes apiserver.
type APIServer struct {
	// URL is the address the ApiServer should listen on for client connections.
	//
	// If this is not specified, we default to a random free port on localhost.
	URL *url.URL

	// Path is the path to the apiserver binary.
	//
	// If this is left as the empty string, we will attempt to locate a binary,
	// by checking for the TEST_ASSET_KUBE_APISERVER environment variable, and
	// the default test assets directory. See the "Binaries" section above (in
	// doc.go) for details.
	Path string

	// ClusterConfig is the kubeadm-compatible configuration for
	// clusters, which is partially supported by this framework.
	//
	// The elements of the ClusterConfig which are supported by
	// this framework are:
	//
	// - ClusterConfig.CertificatesDir
	// - ClusterConfig.APIServerExtraArgs
	ClusterConfig cluster.Config

	// EtcdURL is the URL of the Etcd the APIServer should use.
	//
	// If this is not specified, the Start() method will return an error.
	EtcdURL *url.URL

	// StartTimeout, StopTimeout specify the time the APIServer is allowed to
	// take when starting and stoppping before an error is emitted.
	//
	// If not specified, these default to 20 seconds.
	StartTimeout time.Duration
	StopTimeout  time.Duration

	// Out, Err specify where APIServer should write its StdOut, StdErr to.
	//
	// If not specified, the output will be discarded.
	Out io.Writer
	Err io.Writer

	processState *internal.ProcessState
}

// Start starts the apiserver, waits for it to come up, and returns an error,
// if occurred.
func (s *APIServer) Start() error {
	if s.EtcdURL == nil {
		return fmt.Errorf("expected EtcdURL to be configured")
	}

	var err error

	s.processState = &internal.ProcessState{}

	s.processState.DefaultedProcessInput, err = internal.DoDefaulting(
		"kube-apiserver",
		s.URL,
		s.ClusterConfig.CertificatesDir,
		s.Path,
		s.StartTimeout,
		s.StopTimeout,
	)
	if err != nil {
		return err
	}

	s.processState.HealthCheckEndpoint = "/healthz"

	s.URL = &s.processState.URL
	s.Path = s.processState.Path
	s.StartTimeout = s.processState.StartTimeout
	s.StopTimeout = s.processState.StopTimeout

	tmplData := struct {
		EtcdURL *url.URL
		URL     *url.URL
		CertDir string
	}{
		s.EtcdURL,
		s.URL,
		s.processState.Dir,
	}

	args := flattenArgs(s.ClusterConfig.APIServerExtraArgs)

	s.processState.Args, err = internal.RenderTemplates(
		internal.DoAPIServerArgDefaulting(args), tmplData,
	)
	if err != nil {
		return err
	}

	return s.processState.Start(s.Out, s.Err)
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
