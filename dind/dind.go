package dind

import (
	"bytes"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io"
	"net/url"
	"os"
	"os/exec"

	"sigs.k8s.io/testing_frameworks/cluster"
)

type Dind struct {
	Out io.Writer
	Err io.Writer

	label string
}

func (d *Dind) Setup(c cluster.Config) error {
	label, err := generateRandomString(10)
	if err != nil {
		return err
	}
	d.label = label

	cmd := d.clusterCmd("up")
	d.ensureIO(cmd)
	cmd.Env = d.clusterEnv(
		fmt.Sprintf("NUM_NODES=%d", c.Shape.NodeCount),
	)
	if ok, port := getPortFromURL(c.API.BindURL); ok {
		cmd.Env = append(cmd.Env, fmt.Sprintf("APISERVER_PORT=%s", port))
	}

	return cmd.Run()
}

func (d *Dind) TearDown() error {
	cmd := d.clusterCmd("clean")
	d.ensureIO(cmd)
	cmd.Env = d.clusterEnv()

	return cmd.Run()
}

func (d *Dind) clusterCmd(args ...string) *exec.Cmd {
	binPath := "/Users/pivotal/workspace/kubeadm-dind-cluster/dind-cluster.sh"
	return exec.Command(binPath, args...) // #nosec
}

func (d *Dind) ensureIO(cmd *exec.Cmd) {
	if d.Out != nil {
		cmd.Stdout = d.Out
	}
	if d.Err != nil {
		cmd.Stderr = d.Err
	}
}

func (d *Dind) ClientConfig() *url.URL {
	port, _ := d.getAPIServerPort()
	return &url.URL{
		Scheme: "http",
		Host:   fmt.Sprintf("localhost:%s", port),
	}
}

func (d *Dind) clusterEnv(additionalEnv ...string) []string {
	env := append(os.Environ(),
		fmt.Sprintf("DIND_LABEL=%s", d.label),
	)
	for _, e := range additionalEnv {
		env = append(env, e)
	}
	return env
}

func (d *Dind) getAPIServerPort() (string, error) {
	cmd := d.clusterCmd("apiserver-port")
	cmd.Env = d.clusterEnv()

	stdout := &bytes.Buffer{}
	cmd.Stdout = stdout
	if err := cmd.Run(); err != nil {
		return "", err
	}

	var port int
	_, err := fmt.Fscanf(stdout, "%d\n", &port)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%d", port), nil
}

// ContextName should maybe be implemented in k-d-c
// TODO: move to k-d-c
func (d *Dind) ContextName() string {
	if d.label == "" {
		return "dind"
	}
	hasher := sha1.New()
	hasher.Write([]byte(d.label)) // #nosec, hasher#Write never returns an error
	return fmt.Sprintf("dind-%x", hasher.Sum(nil))
}

func generateRandomString(s int) (string, error) {
	b, err := generateRandomBytes(s)
	return base64.URLEncoding.EncodeToString(b), err
}

func generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func getPortFromURL(u *url.URL) (bool, string) {
	if u == nil {
		return false, ""
	}
	return true, u.Port()
}
