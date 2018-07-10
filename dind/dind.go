package dind

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"os/exec"

	"sigs.k8s.io/testing_frameworks/cluster"
)

type Dind struct {
	port string
	Out  io.Writer
	Err  io.Writer
}

func (d *Dind) Setup(c cluster.Config) error {
	d.port = c.API.BindURL.Port()

	path := "/Users/pivotal/workspace/kubeadm-dind-cluster/dind-cluster.sh"
	cmd := exec.Command(path, "up")
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("APISERVER_PORT=%s", d.port),
		fmt.Sprintf("NUM_NODES=%d", c.Shape.NodeCount),
	)

	if d.Out != nil {
		cmd.Stdout = d.Out
	}
	if d.Err != nil {
		cmd.Stderr = d.Err
	}

	return cmd.Run()
}

func (d *Dind) Teardown() error {
	return nil
}

func (d *Dind) ClientConfig() *url.URL {
	return &url.URL{
		Scheme: "http",
		Host:   fmt.Sprintf("localhost:%s", d.port),
	}
}
