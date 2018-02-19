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

// VirtualKubelet knows how to run a virtual-kubelet
type VirtualKubelet struct {
	// Path is the path to the virtual-kubelet binary.
	//
	// If this is left as the empty string, we will attempt to locate a binary,
	// by checking for the TEST_ASSET_VIRTUAL_KUBELET environment variable, and
	// the default test assets directory. See the "Binaries" section above (in
	// doc.go) for details.
	Path string

	// Args is a list of arguments which will passed to the virtual-kubelet
	// binary.  Before they are passed on, they will be evaluated as go-template
	// strings.
	// This means you can use fields which are defined and exported on this
	// VirtualKubelet struct (e.g. "--kubeconfig={{ .ConfDir }}/kube.conf").
	// Those templates will be evaluated after the defaulting of the
	// VirtualKubelet's fields has already happened and just before the binary
	// actually gets started. Thus you have access to caluclated fields like
	// `ConfDir` and others.
	//
	// If not specified, the minimal set of arguments to run the APIServer will
	// be used.
	Args []string

	// ConfDir is a path to a directory containing whatever configuration the
	// virtual kubelet needs.
	//
	// If left unspecified, then the Start() method will create a fresh temporary
	// directory, and the Stop() method will clean it up.
	ConfDir string

	// StartTimeout, StopTimeout specify the time the VirtualKubelet is allowed
	// to take when starting and stopping before an error is emitted.
	//
	// If not specified, these default to 20 seconds.
	StartTimeout time.Duration
	StopTimeout  time.Duration

	// Out, Err specify where VirtualKubelet should write its StdOut, StdErr to.
	//
	// If not specified, the output will be discarded.
	Out io.Writer
	Err io.Writer

	processState *internal.ProcessState
}

// Start starts the virtual kubelet, waits for it to come up, and returns an
// error, if one occoured.
func (vk *VirtualKubelet) Start(apiServerConnectionConfig RemoteConnectionConfig) error {
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

	templateVars := struct {
		*internal.ProcessState
	}{
		vk.processState,
	}

	vk.processState.Args, err = internal.RenderTemplates(
		doArgDefaulting(vk.Args, VirtualKubeletDefaultArgs),
		templateVars,
	)
	if err != nil {
		return err
	}

	confPath := path.Join(vk.processState.Dir, "kube.conf")
	if err := writeCubeConfig(confPath, apiServerConnectionConfig.URL); err != nil {
		return err
	}

	return vk.processState.Start(vk.Out, vk.Err)
}

// Stop stops this process gracefully, waits for its termination, and cleans up
// the ConfDir if necessary.
func (vk *VirtualKubelet) Stop() error {
	return vk.processState.Stop()
}

// VirtualKubeletDefaultArgs is the default set of arguments that get passed
// to the VirtualKubelet binary.
var VirtualKubeletDefaultArgs = []string{
	"--provider=mock",
	"--kubeconfig={{ .Dir }}/kube.conf",
}

func writeCubeConfig(kubeConfPath string, apiServerURL *url.URL) error {
	var err error
	var file *os.File

	file, err = os.Create(kubeConfPath)
	if err != nil {
		return err
	}

	kubeConf := fmt.Sprintf(kubeConfTmpl, apiServerURL)
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
