package integration_tests

import (
	"bytes"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"sigs.k8s.io/testing_frameworks/cluster"
	. "sigs.k8s.io/testing_frameworks/lightweight"
)

var _ = Describe("Etcd", func() {
	It("can inspect IO", func() {
		stderr := &bytes.Buffer{}

		config := cluster.Config{}
		config.Etcd.ProcessConfig.Err = stderr

		etcd := &Etcd{
			ClusterConfig: config,
		}

		Expect(etcd.Start()).To(Succeed())
		defer func() {
			Expect(etcd.Stop()).To(Succeed())
		}()

		Expect(stderr.String()).NotTo(BeEmpty())
	})

	It("can use user specified Args", func() {
		stdout := &bytes.Buffer{}
		stderr := &bytes.Buffer{}

		config := cluster.Config{}
		config.Etcd.ExtraArgs = map[string]string{
			"--help": "",
		}
		config.Etcd.ProcessConfig.Out = stdout
		config.Etcd.ProcessConfig.Err = stderr
		config.Etcd.ProcessConfig.StartTimeout = 500 * time.Millisecond

		etcd := &Etcd{
			ClusterConfig: config,
		}

		// it will timeout, as we'll never see the "startup message" we are waiting
		// for on StdErr
		Expect(etcd.Start()).To(MatchError(ContainSubstring("timeout")))

		Expect(stdout.String()).To(ContainSubstring("member flags"))
		Expect(stderr.String()).To(ContainSubstring("usage: etcd"))
	})
})
