/*
Copyright 2018 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package internal

import (
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

var rootPath string

func init() {
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		panic("Could not determine the path of the BinPathFinder")
	}
	rootPath = filepath.Join(filepath.Dir(thisFile), "..")
}

// BinPathFinder derives an environment variable based on the symbolic name;
// if the environment variable is set, it uses its value as the binary path
// if it's not set, it falls back to a default assets location in the
// <containingDirectory>
func BinPathFinder(containingDirectory, symbolicName string) (binPath string) {
	punctuationPattern := regexp.MustCompile("[^A-Z0-9]+")
	sanitizedName := punctuationPattern.ReplaceAllString(strings.ToUpper(symbolicName), "_")
	leadingNumberPattern := regexp.MustCompile("^[0-9]+")
	sanitizedName = leadingNumberPattern.ReplaceAllString(sanitizedName, "")
	envVar := "TEST_ASSET_" + sanitizedName

	if val, ok := os.LookupEnv(envVar); ok {
		return val
	}

	return filepath.Join(rootPath, containingDirectory, "assets", "bin", symbolicName)
}
