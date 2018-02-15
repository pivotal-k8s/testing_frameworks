package integration

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"path"
	"time"

	"github.com/kubernetes-sig-testing/frameworks/integration/internal"
)

type VirtualKubelet struct {
	Path         string
	Args         []string
	ConfDir      string
	APIServerURL *url.URL
	StartTimeout time.Duration
	StopTimeout  time.Duration
	Out          io.Writer
	Err          io.Writer

	processState *internal.ProcessState
}

func (vk *VirtualKubelet) Start() error {
	var err error

	vk.processState = &internal.ProcessState{}

	vk.processState.DefaultedProcessInput, err = internal.DoDefaulting(
		"virtual-kubelet",
		nil,
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

func (vk *VirtualKubelet) RegisterTo(cp *ControlPlane) {
	vk.APIServerURL = cp.APIURL()
}

func (vk *VirtualKubelet) setConf() error {
	kubeConfPath := path.Join(vk.ConfDir, "kube.conf")

	var err error
	var file *os.File

	file, err = os.Create(kubeConfPath)
	if err != nil {
		return err
	}

	kubeConf := fmt.Sprintf(kubeConfTmpl, vk.APIServerURL)
	_, err = file.Write([]byte(kubeConf))
	if err != nil {
		return err
	}

	err = file.Close()
	if err != nil {
		return err
	}

	return nil
}

const kubeConfTmpl string = `
apiVersion: v1
kind: Config
users:
- name: vk_user
  user:
    username: admin
    password: admin
clusters:
- name: vk_cluster
  cluster:
    server: %s
contexts:
- context:
    cluster: vk_cluster
    user: vk_user
  name: vk_ctx
current-context: vk_ctx
`
