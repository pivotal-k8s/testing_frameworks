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

readonly KUBEADM_VERSION="${KUBEADM_VERSION:=v1alpha3}"
readonly KUBEADM_BASE_URL="https://raw.githubusercontent.com/kubernetes/kubernetes/master/cmd/kubeadm/app/apis/kubeadm/${KUBEADM_VERSION}"
readonly CLIENTCMD_VERSION="${CLIENTCMD_VERSION:=v1}"
readonly CLIENTCMD_BASE_URL="https://raw.githubusercontent.com/kubernetes/client-go/master/tools/clientcmd/api/${CLIENTCMD_VERSION}"
readonly DEST="$( dirname "$0" )"
readonly PKG="$( basename "$DEST" )"

dl() {
  curl -sL "$1"
}

getKubeadmReplacements() {
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

getClientcmdReplacements() {
  cat <<'EOF'
s@\(^[^/].*\bruntime\b.*\)@// TODO \1@g
EOF
}

replaceUnsupported() {
  sed -f "$1"
}

setPackage() {
  sed "s/^package .*/package ${1}/g"
}

transform() {
  replaceUnsupported "$1" \
    | setPackage "$2" \
    | gofmt -s \
    | goimports
}

main() {
  dl "${CLIENTCMD_BASE_URL}/types.go" \
    | transform <(getClientcmdReplacements) "$PKG" \
    > "${DEST}/clientcmd.go"

  dl "${KUBEADM_BASE_URL}/types.go" \
    | transform <(getKubeadmReplacements) "$PKG" \
    > "${DEST}/kubeadm.go"
}

main "$@"
