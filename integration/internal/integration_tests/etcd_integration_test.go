package integration_tests

import (
	"bytes"
	"io/ioutil"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/kubernetes-sig-testing/frameworks/integration"
)

var _ = Describe("Etcd", func() {
	It("uses the configured datadir", func() {
		tmpDir, err := ioutil.TempDir("", "k8s_test_framework_")
		Expect(err).NotTo(HaveOccurred())

		etcd := &Etcd{}
		etcd.DataDir = tmpDir

		Expect(etcd.Start()).To(Succeed())
		defer func() {
			Expect(etcd.Stop()).To(Succeed())
		}()

		Expect(directoryHasContents(tmpDir)).To(BeTrue())
	})

	It("can inspect IO", func() {
		stderr := &bytes.Buffer{}
		etcd := &Etcd{
			Err: stderr,
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
		etcd := &Etcd{
			Args:         []string{"--help"},
			Out:          stdout,
			Err:          stderr,
			StartTimeout: 500 * time.Millisecond,
		}

		// it will timeout, as we'll never see the "startup message" we are waiting
		// for on StdErr
		Expect(etcd.Start()).To(MatchError(ContainSubstring("timeout")))

		Expect(stdout.String()).To(ContainSubstring("member flags"))
		Expect(stderr.String()).To(ContainSubstring("usage: etcd"))
	})
})

func directoryHasContents(dir string) bool {
	fileList, err := ioutil.ReadDir(dir)
	Expect(err).NotTo(HaveOccurred(), "Cannot read directory %s", dir)
	return len(fileList) != 0
}
