package test

import (
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

var assetsPath string

func init() {
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		panic("Could not determine the path of the BinPathFinder")
	}
	assetsPath = filepath.Join(filepath.Dir(thisFile), "assets", "bin")
}

// BinPathFinder is the signature of a function translating a symbolic name of a binary to
// a resolvable path to the binary in question
type BinPathFinder func(symbolicName string) (binPath string)

// DefaultBinPathFinder is an implementation of BinPathFinder which checks the an environment
// variable, derived from the symbolic name, and falls back to a default assets location when
// this variable is not set
func DefaultBinPathFinder(symbolicName string) (binPath string) {
	punctuationPattern := regexp.MustCompile("[^A-Z0-9]+")
	sanitizedName := punctuationPattern.ReplaceAllString(strings.ToUpper(symbolicName), "_")
	leadingNumberPattern := regexp.MustCompile("^[0-9]+")
	sanitizedName = leadingNumberPattern.ReplaceAllString(sanitizedName, "")
	envVar := "TEST_ASSET_" + sanitizedName

	if val, ok := os.LookupEnv(envVar); ok {
		return val
	}

	return filepath.Join(assetsPath, symbolicName)
}
