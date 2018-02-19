package integration

import (
	"errors"
	"net/url"

	"github.com/kubernetes-sig-testing/frameworks/integration/internal"
)

// RemoteConnectionConfig is a struct holding certain configuration items
// describing how to connect to a remote endpoint.
type RemoteConnectionConfig struct {
	URL *url.URL
}

func processStateToConnectionConfig(ps *internal.ProcessState) (RemoteConnectionConfig, error) {
	if ps == nil {
		return RemoteConnectionConfig{}, errors.New("no process state; did you call Start()?")
	}
	if ps.URL == (url.URL{}) {
		return RemoteConnectionConfig{}, errors.New("Process has not bound to an URL yet; did you call Start()?")
	}

	return RemoteConnectionConfig{URL: &ps.URL}, nil
}
