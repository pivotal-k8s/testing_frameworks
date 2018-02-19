package integration

import (
	"fmt"
	"net/url"

	"github.com/kubernetes-sig-testing/frameworks/integration/internal"
)

// RemoteConnectionConfig is a struct holding certain configuration items
// describing how to connect to a remote endpoint.
type RemoteConnectionConfig struct {
	URL *url.URL
}

func processStateToConnectionConfig(ps *internal.ProcessState) (RemoteConnectionConfig, error) {
	if ps.URL == (url.URL{}) {
		return RemoteConnectionConfig{}, fmt.Errorf("Process has not bound to an URL yet; did you call Start()?")
	}

	return RemoteConnectionConfig{URL: &ps.URL}, nil
}
