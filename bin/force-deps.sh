#!/bin/bash

set -eu

[ -z "${DEBUG:-}" ] || {
  set -x
}


##------------------------------------------------------------
BASE_DIR="$( cd "$(dirname "$0")/.." && pwd )"
PATH="${PATH}:${GOPATH}/bin"

##------------------------------------------------------------
forceClientGoVersion() {
  local package="k8s.io/client-go"
  local version="${1:-}"

  rm -rf "${BASE_DIR}/vendor/${package}"

  go get "${package}/..."
  (
    cd "${GOPATH}/src/${package}"
    git checkout "$version" >/dev/null 2>&1
    godep restore ./...
  )
}

##------------------------------------------------------------
main() {
  [ -z "${CLIENT_GO_VERSION:-}" ] && {
    echo 'Using client-go from ./vendor' >&2
    exit
  }

  echo "Setting client-go version to ${CLIENT_GO_VERSION}" >&2
  forceClientGoVersion "$CLIENT_GO_VERSION"
}

##------------------------------------------------------------
##------------------------------------------------------------
main "$@"
