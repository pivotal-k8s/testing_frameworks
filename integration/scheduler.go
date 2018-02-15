package integration

import (
	"fmt"
	"io"
	"net/url"
	"time"

	"github.com/kubernetes-sig-testing/frameworks/integration/internal"
)

// Scheduler is kube-scheduler. It can be used as an `AdditionalComponent` for
// the control plane.
type Scheduler struct {
	// URL is the address the Scheduler should listen on for connections.
	//
	// If this is not specified, we default to a random free port on localhost.
	URL *url.URL

	// Path is the path to the scheduler binary.
	//
	// If this is left as the empty string, we will attempt to locate a binary,
	// by checking for the TEST_ASSET_KUBE_SCHEDULER environment variable, and
	// the default test assets directory. See the "Binaries" section above (in
	// doc.go) for details.
	Path string

	// Args is a list of arguments which will passed to the Scheduler binary.
	// Before they are passed on, they will be evaluated as go-template strings.
	// This means you can use fields which are defined and exported on this
	// Scheduler struct (e.g. "--master={{ .APIServerURL.String }}",
	// Those templates will be evaluated after the defaulting of the Scheduler's
	// fields has already happened and just before the binary actually gets
	// started. Thus you have access to caluclated fields like `URL` and others.
	//
	// If not specified, the minimal set of arguments to run the Scheduler will
	// be used.
	Args []string

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

// Start starts the scheduler, waits for it to come up, and returns an error,
// if occurred.
func (c *Scheduler) Start(r RemoteConnectionConfig) error {
	var err error

	if r.URL == nil {
		return fmt.Errorf("Remote connection config must include a URL")
	}

	c.processState = &internal.ProcessState{}

	c.processState.DefaultedProcessInput, err = internal.DoDefaulting(
		"kube-scheduler",
		c.URL,
		"",
		false,
		c.Path,
		c.StartTimeout,
		c.StopTimeout,
	)
	if err != nil {
		return err
	}

	c.processState.StartMessage = "starting healthz server on"

	// TODO Remove mutation of main struct
	c.URL = &c.processState.URL

	templateVars := struct {
		*internal.ProcessState
		APIServerURL *url.URL
	}{
		c.processState,
		r.URL,
	}

	c.processState.Args, err = internal.RenderTemplates(
		internal.DoSchedulerArgDefaulting(c.Args),
		templateVars,
	)
	if err != nil {
		return err
	}

	return c.processState.Start(c.Out, c.Err)
}

// Stop stops this process gracefully, waits for its termination.
func (c *Scheduler) Stop() error {
	return c.processState.Stop()
}
