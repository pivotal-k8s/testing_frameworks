package integration

import (
	"io"
	"net/url"
	"os"
	"path"
	"time"

	"github.com/kubernetes-sig-testing/frameworks/integration/internal"
)

type VirtualKubelet struct {
	URL          *url.URL
	Path         string
	Args         []string
	ConfDir      string
	StartTimeout time.Duration
	StopTimeout  time.Duration
	Out          io.Writer
	Err          io.Writer

	Conf string

	processState *internal.ProcessState
}

func (vk *VirtualKubelet) Start() error {
	var err error

	vk.processState = &internal.ProcessState{}

	vk.processState.DefaultedProcessInput, err = internal.DoDefaulting(
		"virtual-kubelet",
		vk.URL,
		vk.ConfDir,
		true,
		vk.Path,
		vk.StartTimeout,
		vk.StopTimeout,
	)
	if err != nil {
		return err
	}

	vk.processState.StartMessage = "Node 'virtual-kubelet' with OS type 'Linux' registered"

	vk.URL = &vk.processState.URL
	vk.Path = vk.processState.Path
	vk.StartTimeout = vk.processState.StartTimeout
	vk.StopTimeout = vk.processState.StopTimeout
	vk.ConfDir = vk.processState.Dir

	vk.processState.Args, err = internal.RenderTemplates(
		internal.DoVirtualKubeletArgDefaulting(vk.Args), vk,
	)
	if err != nil {
		return err
	}

	if err := vk.setConf(); err != nil {
		return err
	}

	return vk.processState.Start(vk.Out, vk.Err)
}

func (vk *VirtualKubelet) Stop() error {
	return vk.processState.Stop()
}

func (vk *VirtualKubelet) setConf() error {
	kubeConfPath := path.Join(vk.ConfDir, "kube.conf")

	var err error
	var file *os.File

	file, err = os.Create(kubeConfPath)
	if err != nil {
		return err
	}

	_, err = file.Write([]byte(vk.Conf))
	if err != nil {
		return err
	}

	err = file.Close()
	if err != nil {
		return err
	}

	return nil
}
