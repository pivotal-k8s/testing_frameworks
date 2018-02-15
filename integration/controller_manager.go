package integration

import (
	"fmt"
	"io"
	"net/url"
	"time"

	"github.com/kubernetes-sig-testing/frameworks/integration/internal"
)

// ControllerManager is a kube-controller-manager. It can be used as an
// `AdditionalComponent` for the control plane.
type ControllerManager struct {
	// URL is the address the ControllerManager should listen on for connections.
	//
	// If this is not specified, we default to a random free port on localhost.
	URL *url.URL

	// Path is the path to the ControllerManager binary.
	//
	// If this is left as the empty string, we will attempt to locate a binary,
	// by checking for the TEST_ASSET_KUBE_CONTROLLER_MANAGER environment
	// variable, and the default test assets directory. See the "Binaries"
	// section above (in doc.go) for details.
	Path string

	// Args is a list of arguments which will passed to the ControllerManager
	// binary. Before they are passed on, they will be evaluated as go-template
	// strings.
	// This means you can use fields which are defined and exported on this
	// ControllerManager struct (e.g. "--master={{ .APIServerURL.String }}",
	// Those templates will be evaluated after the defaulting of the
	// ControllerManager's fields has already happened and just before the binary
	// actually gets started. Thus you have access to caluclated fields like
	// `URL` and others.
	//
	// If not specified, the minimal set of arguments to run the
	// ControllerManager will be used.
	Args []string

	// APIServerURL is the URL pointing to the APIServer of the control plane
	// this Scheduler should connect to.
	//
	// It can be set directly on the struct, or indirectly by calling the
	// `RegisterTo` method.
	// This setting is mandatory and an error will be emitted when this is not
	// set at start time.
	APIServerURL *url.URL

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

// Start starts the controller manager, waits for it to come up, and returns an
// error, if occurred.
func (c *ControllerManager) Start() error {
	var err error

	if c.APIServerURL == nil {
		return fmt.Errorf("APIServerURL must be configured")
	}

	c.processState = &internal.ProcessState{}

	c.processState.DefaultedProcessInput, err = internal.DoDefaulting(
		"kube-controller-manager",
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

	c.processState.StartMessage = "Sending events to api server."

	c.URL = &c.processState.URL
	c.Path = c.processState.Path
	c.StartTimeout = c.processState.StartTimeout
	c.StopTimeout = c.processState.StopTimeout

	c.processState.Args, err = internal.RenderTemplates(
		internal.DoControllerManagerArgDefaulting(c.Args), c,
	)
	if err != nil {
		return err
	}

	return c.processState.Start(c.Out, c.Err)
}

// Stop stops this process gracefully, waits for its termination.
func (c *ControllerManager) Stop() error {
	return c.processState.Stop()
}

// RegisterTo configures the controller manager in a way so that it registers &
// connects to the control plane handed in upon start. It is the responsibility
// of this controller manager to get all the data needed from the control plane
// and set itself up accordingly.
func (c *ControllerManager) RegisterTo(cp *ControlPlane) {
	c.APIServerURL = cp.APIURL()
}
