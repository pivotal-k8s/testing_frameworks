#!/usr/bin/env bash

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
    > "${DEST}/kubeadm.go"
}

main "$@"
