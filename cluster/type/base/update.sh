#!/usr/bin/env bash

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

set -e
set -u
set -o pipefail

[ -n "${DEBUG:-}" ] && set -x

readonly VERSION="${VERSION:=}"
readonly BASE_URL="https://raw.githubusercontent.com/kubernetes/kubernetes/master/cmd/kubeadm/app/apis/kubeadm/${VERSION}"
readonly DEST="$( dirname "$0" )"
readonly PKG="$( basename "$DEST" )"

dl() {
  curl -sL "$1"
}

getReplacements() {
  cat <<'EOF'

s@\(^[^/].*\bfuzz\b.*\)@// TODO \1@g
s@\(^[^/].*\bmetav1\b.*\)@// TODO \1@g
s@\(^[^/].*\bv1\b.*\)@// TODO \1@g
s@\(^[^/].*\bkubeletconfigv1beta1\b.*\)@// TODO \1@g
s@\(^[^/].*\bkubeletconfig\b.*\)@// TODO \1@g
s@\(^[^/].*\bkubeproxyconfigv1alpha1\b.*\)@// TODO \1@g
s@\(^[^/].*\bkubeproxyconfig\b.*\)@// TODO \1@g

EOF
}

cleanMeta() {
  sed -f <( getReplacements )
}

setPackage() {
  sed "s/^package .*/package ${PKG}/g"
}

main() {
  dl "${BASE_URL}/types.go" \
    | cleanMeta \
    | setPackage \
    | gofmt -s \
    | goimports \
    > "${DEST}/kubeadm.go"
}

main "$@"
