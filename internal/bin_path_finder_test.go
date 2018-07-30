package internal

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("BinPathFinder", func() {
	Context("when relying on the default assets path", func() {
		var (
			previousAssetsPath string
		)
		BeforeEach(func() {
			previousAssetsPath = rootPath
			rootPath = "/some/path"
		})
		AfterEach(func() {
			rootPath = previousAssetsPath
		})
		It("returns the default path when no env var is configured", func() {
			binPath := BinPathFinder("some_part", "some_bin")
			Expect(binPath).To(Equal("/some/path/some_part/assets/bin/some_bin"))
		})
	})

	Context("when environment is configured", func() {
		var (
			previousValue string
			wasSet        bool
		)
		BeforeEach(func() {
			envVarName := "TEST_ASSET_ANOTHER_SYMBOLIC_NAME"
			if val, ok := os.LookupEnv(envVarName); ok {
				previousValue = val
				wasSet = true
			}
			os.Setenv(envVarName, "/path/to/some_bin.exe")
		})
		AfterEach(func() {
			if wasSet {
				os.Setenv("TEST_ASSET_ANOTHER_SYMBOLIC_NAME", previousValue)
			} else {
				os.Unsetenv("TEST_ASSET_ANOTHER_SYMBOLIC_NAME")
			}
		})
		It("returns the path from the env", func() {
			binPath := BinPathFinder("does_not_matter", "another_symbolic_name")
			Expect(binPath).To(Equal("/path/to/some_bin.exe"))
		})

		It("sanitizes the environment variable name", func() {
			By("cleaning all non-underscore punctuation")
			binPath := BinPathFinder("does_not_matter", "another-symbolic name")
			Expect(binPath).To(Equal("/path/to/some_bin.exe"))
			binPath = BinPathFinder("does_not_matter", "another+symbolic\\name")
			Expect(binPath).To(Equal("/path/to/some_bin.exe"))
			binPath = BinPathFinder("does_not_matter", "another=symbolic.name")
			Expect(binPath).To(Equal("/path/to/some_bin.exe"))
			By("removing numbers from the beginning of the name")
			binPath = BinPathFinder("does_not_matter", "12another_symbolic_name")
			Expect(binPath).To(Equal("/path/to/some_bin.exe"))
		})
	})
})
