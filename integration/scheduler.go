package integration

import (
	"fmt"
	"io"
	"net/url"
	"time"

	"github.com/kubernetes-sig-testing/frameworks/integration/internal"
)

type Scheduler struct {
	URL          *url.URL
	Path         string
	Args         []string
	APIServerURL *url.URL
	StartTimeout time.Duration
	StopTimeout  time.Duration
	Out          io.Writer
	Err          io.Writer

	processState *internal.ProcessState
}

func (c *Scheduler) Start() error {
	var err error

	if c.APIServerURL == nil {
		return fmt.Errorf("APIServerURL must be configured")
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

	c.URL = &c.processState.URL
	c.Path = c.processState.Path
	c.StartTimeout = c.processState.StartTimeout
	c.StopTimeout = c.processState.StopTimeout

	c.processState.Args, err = internal.RenderTemplates(
		internal.DoSchedulerArgDefaulting(c.Args), c,
	)
	if err != nil {
		return err
	}

	return c.processState.Start(c.Out, c.Err)
}

func (c *Scheduler) Stop() error {
	return c.processState.Stop()
}

func (c *Scheduler) RegisterTo(cp *ControlPlane) {
	c.APIServerURL = cp.APIURL()
}
