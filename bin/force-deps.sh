#!/bin/bash

set -eu

[ -z "${DEBUG:-}" ] || {
  set -x
}


##------------------------------------------------------------
BASE_DIR="$( cd "$(dirname "$0")/.." && pwd )"
PATH="${PATH}:${GOPATH}/bin"

##------------------------------------------------------------
main() {
  [ -z "${CLIENT_GO_VERSION:-}" ] || {
    {
      echo "Setting client-go version to ${CLIENT_GO_VERSION}"
      echo
    }>&2

    "${BASE_DIR}/bin/change-client-go-version.sh" \
      "$CLIENT_GO_VERSION" \
      "${BASE_DIR}/Gopkg.toml"
  }

  dep status
}

##------------------------------------------------------------
##------------------------------------------------------------
main "$@"
