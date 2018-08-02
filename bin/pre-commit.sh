#!/bin/bash

# Copyright 2018 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Make sure, we run in the root of the repo and
# therefore run the tests on all packages
base_dir="$( cd "$(dirname "$0")/.." && pwd )"
cd "$base_dir" || {
  echo "Cannot cd to '$base_dir'. Aborting." >&2
  exit 1
}

rc=0

go_dirs() {
  go list -f '{{.Dir}}' ./... | tr '\n' '\0'
}

echo "Update submodules"
git submodule update --init --recursive
rc=$((rc || $?))

echo "Verify Boilerplate"
"${base_dir}/bin/verify-boilerplate.sh" --rootdir="$base_dir" --boilerplate-dir="${base_dir}/bin/boilerplate/"
rc=$((rc || $?))

echo "Running go fmt"
diff <(echo -n) <(go_dirs | xargs -0 gofmt -s -d -l)
rc=$((rc || $?))

echo "Running goimports"
diff -u <(echo -n) <(go_dirs | xargs -0 goimports -l)
rc=$((rc || $?))

echo "Running go vet"
go vet -all ./...
rc=$((rc || $?))

echo "Installing test binaries"
./lightweight/scripts/download-binaries.sh
rc=$((rc || $?))

echo "Running go test"
go test -v ./...
rc=$((rc || $?))

exit $rc
