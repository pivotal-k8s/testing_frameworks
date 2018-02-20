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
			previousAssetsPath = defaultAssetsPath
			defaultAssetsPath = "/some/path/assets/bin"
		})
		AfterEach(func() {
			defaultAssetsPath = previousAssetsPath
		})
		It("returns the default path when no env var is configured", func() {
			binPath := BinPathFinder("some_bin")
			Expect(binPath).To(Equal("/some/path/assets/bin/some_bin"))
		})
	})

	Context("when a custom assets directory is configured", func() {
		BeforeEach(func() {
			os.Setenv("TEST_ASSETS_PATH", "/custom/assets/dir")
		})

		AfterEach(func() {
			os.Unsetenv("TEST_ASSETS_PATH")
		})

		It("returns the path to a binary in the custom assets directory", func() {
			binPath := BinPathFinder("some_bin")
			Expect(binPath).To(Equal("/custom/assets/dir/some_bin"))
		})

		Context("and a custom binary path is configured", func() {
			BeforeEach(func() {
				os.Setenv("TEST_ASSET_SOME_BIN", "/custom/path/to/some_bin")
			})

			AfterEach(func() {
				os.Unsetenv("TEST_ASSET_SOME_BIN")
			})

			It("returns the custom binary path", func() {
				binPath := BinPathFinder("some_bin")
				Expect(binPath).To(Equal("/custom/path/to/some_bin"))
			})
		})
	})

	Context("when a binary environment variable is configured", func() {
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
			binPath := BinPathFinder("another_symbolic_name")
			Expect(binPath).To(Equal("/path/to/some_bin.exe"))
		})

		It("sanitizes the environment variable name", func() {
			By("cleaning all non-underscore punctuation")
			binPath := BinPathFinder("another-symbolic name")
			Expect(binPath).To(Equal("/path/to/some_bin.exe"))
			binPath = BinPathFinder("another+symbolic\\name")
			Expect(binPath).To(Equal("/path/to/some_bin.exe"))
			binPath = BinPathFinder("another=symbolic.name")
			Expect(binPath).To(Equal("/path/to/some_bin.exe"))
			By("removing numbers from the beginning of the name")
			binPath = BinPathFinder("12another_symbolic_name")
			Expect(binPath).To(Equal("/path/to/some_bin.exe"))
		})
	})
})
