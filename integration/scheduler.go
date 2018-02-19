package integration

import (
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

	// StartTimeout, StopTimeout specify the time the Scheduler is allowed to
	// take when starting and stopping before an error is emitted.
	//
	// If not specified, these default to 20 seconds.
	StartTimeout time.Duration
	StopTimeout  time.Duration

	// Out, Err specify where Scheduler should write its StdOut, StdErr to.
	//
	// If not specified, the output will be discarded.
	Out io.Writer
	Err io.Writer

	processState *internal.ProcessState
}

// Start starts the scheduler, waits for it to come up, and returns an error,
// if occurred.
func (s *Scheduler) Start(apiServerConnectionConfig RemoteConnectionConfig) error {
	var err error

	s.processState = &internal.ProcessState{}

	s.processState.DefaultedProcessInput, err = internal.DoDefaulting(
		"kube-scheduler",
		s.URL,
		"",
		false,
		s.Path,
		s.StartTimeout,
		s.StopTimeout,
	)
	if err != nil {
		return err
	}

	s.processState.StartMessage = "starting healthz server on"

	templateVars := struct {
		*internal.ProcessState
		APIServerURL *url.URL
	}{
		s.processState,
		apiServerConnectionConfig.URL,
	}

	s.processState.Args, err = internal.RenderTemplates(
		s.doArgDefaulting(),
		templateVars,
	)
	if err != nil {
		return err
	}

	return s.processState.Start(s.Out, s.Err)
}

// Stop stops this process gracefully, waits for its termination.
func (s *Scheduler) Stop() error {
	return s.processState.Stop()
}

// ConnectionConfig returns the configuration needed to connect to this
// Scheduler.
func (s *Scheduler) ConnectionConfig() (conf RemoteConnectionConfig, err error) {
	return processStateToConnectionConfig(s.processState)
}

func (s *Scheduler) doArgDefaulting() []string {
	if len(s.Args) != 0 {
		return s.Args
	}

	return SchedulerDefaultArgs
}

// SchedulerDefaultArgs is the default set of arguments that get passed to the
// Scheduler binary.
var SchedulerDefaultArgs = []string{
	"--master={{ .APIServerURL.String }}",
	"--port={{ .URL.Port }}",
	"--address={{ .URL.Hostname }}",
}
