package integration

import (
	"fmt"
	"io"
	"time"

	"net/url"

	"github.com/kubernetes-sig-testing/frameworks/integration/internal"
)

// Etcd knows how to run an etcd server.
type Etcd struct {
	// URL is the address the Etcd should listen on for client connections.
	//
	// If this is not specified, we default to a random free port on localhost.
	URL *url.URL

	// Path is the path to the etcd binary.
	//
	// If this is left as the empty string, we will attempt to locate a binary,
	// by checking for the TEST_ASSET_ETCD environment variable, and the default
	// test assets directory. See the "Binaries" section above (in doc.go) for
	// details.
	Path string

	// DataDir is a path to a directory in which etcd can store its state.
	//
	// If left unspecified, then the Start() method will create a fresh temporary
	// directory, and the Stop() method will clean it up.
	DataDir string

	// StartTimeout, StopTimeout specify the time the Etcd is allowed to
	// take when starting and stopping before an error is emitted.
	//
	// If not specified, these default to 20 seconds.
	StartTimeout time.Duration
	StopTimeout  time.Duration

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
		e.DataDir,
		e.Path,
		e.StartTimeout,
		e.StopTimeout,
	)
	if err != nil {
		return err
	}

	e.processState.Args = internal.MakeEtcdArgs(e.processState.DefaultedProcessInput)

	e.processState.StartMessage = fmt.Sprintf("serving insecure client requests on %s", e.processState.URL.Hostname())

	return e.processState.Start()
}

// Stop stops this process gracefully, waits for its termination, and cleans up
// the DataDir if necessary.
func (e *Etcd) Stop() error {
	return e.processState.Stop()
}

// ListeningURL returns the URL this Etcd actually listens on.
//
// If you specified
// a URL at construction time, that's what will be returned here. If Etcd has
// defaulted to a free port on localhost, that URL is returned. If no URL has
// been specified, and Etcd hasn't started yet, then we return an error.
func (e *Etcd) ListeningURL() (*url.URL, error) {
	return e.processState.ListeningURL(e.URL)
}

// StdOut returns a Reader for the StdOut stream.
//
// If you haven't started your Etcd yet, then this method will return an error.
func (e *Etcd) StdOut() (io.Reader, error) {
	return e.processState.StdOut()
}

// StdErr returns a Reader for the StdErr stream.
//
// If you haven't started your Etcd yet, then this method will return an error.
func (e *Etcd) StdErr() (io.Reader, error) {
	return e.processState.StdErr()
}
