package integration

import (
	"io"
	"net/url"
	"time"

	"github.com/kubernetes-sig-testing/frameworks/integration/internal"
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

	EnableRBAC bool

	// Args is a list of arguments which will passed to the APIServer binary.
	// Before they are passed on, they will be evaluated as go-template strings.
	// This means you can use fields which are defined and exported on this
	// APIServer struct (e.g. "--cert-dir={{ .Dir }}").
	// Those templates will be evaluated after the defaulting of the APIServer's
	// fields has already happened and just before the binary actually gets
	// started. Thus you have access to caluclated fields like `URL` and others.
	//
	// If not specified, the minimal set of arguments to run the APIServer will
	// be used.
	Args []string

	// CertDir is a path to a directory containing whatever certificates the
	// APIServer will need.
	//
	// If left unspecified, then the Start() method will create a fresh temporary
	// directory, and the Stop() method will clean it up.
	CertDir string

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
func (s *APIServer) Start(etcdConnectionConfig RemoteConnectionConfig) error {
	var err error

	s.processState = &internal.ProcessState{}

	s.processState.DefaultedProcessInput, err = internal.DoDefaulting(
		"kube-apiserver",
		s.URL,
		s.CertDir,
		true,
		s.Path,
		s.StartTimeout,
		s.StopTimeout,
	)
	if err != nil {
		return err
	}

	s.processState.StartMessage = internal.GetAPIServerStartMessage(s.processState.URL)

	templateVars := struct {
		EtcdURL *url.URL
		*internal.ProcessState
	}{
		etcdConnectionConfig.URL,
		s.processState,
	}

	s.processState.Args, err = internal.RenderTemplates(
		internal.DoAPIServerArgDefaulting(s.Args),
		templateVars,
	)
	if err != nil {
		return err
	}
	if s.EnableRBAC == true {
		s.processState.Args = append(
			s.processState.Args,
			"--authorization-mode=RBAC",
		)
	}

	return s.processState.Start(s.Out, s.Err)
}

// Stop stops this process gracefully, waits for its termination, and cleans up
// the CertDir if necessary.
func (s *APIServer) Stop() error {
	return s.processState.Stop()
}

func (s *APIServer) ConnectionConfig() (RemoteConnectionConfig, error) {
	return processStateToConnectionConfig(s.processState)
}
