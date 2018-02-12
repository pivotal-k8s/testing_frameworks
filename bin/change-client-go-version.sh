#!/usr/bin/env bash

set -eu
set -o pipefail

[ -z "${DEBUG:-}" ] || {
  set -x
}

readonly GOPKG_TMPL='
[[override]]
  name = "%s"
  revision = "%s"
'

printDepsByRev() {
  jq -r '.Deps[] | "\(.Rev) \(.ImportPath)"' | sort
}

combineByRev() {
  awk -vTMPL="$GOPKG_TMPL" '
    {
      if (tuples[$1]) {
        tuples[$1] = tuples[$1] "\0" $2
      } else {
        tuples[$1] = $2
      }
    }

    END {
      for (rev in tuples) {
        package = pkg2repo(lcp(tuples[rev]))
        printf TMPL, package, rev
      }
    }

    function pkg2repo(package,  repo,len,tokens) {
      gsub(/\/$/, "", package)
      len = split(package, tokens, "/")

      if (tokens[1] == "github.com") {
        return tokens[1] "/" tokens[2] "/" tokens[3]
      }

      if (len == 1) {
        return tokens[1]
      }

      if (len == 2) {
        return tokens[1] "/" tokens[2]
      }

      if (tokens[3] == "pkg") {
        return tokens[1] "/" tokens[2]
      }

      return tokens[1] "/" tokens[2] "/" tokens[3]
    }

    # https://rosettacode.org/wiki/Longest_common_prefix#AWK
    function lcp(str,  arr,hits,i,j,lcp_leng,n,sw_leng) {
      n = split(str,arr,"\0")
      if (n == 0) { # null string
        return("")
      }
      if (n == 1) { # only 1 word, then it is the longest
        return(str)
      }
      sw_leng = length(arr[1])
      for (i=2; i<=n; i++) { # find shortest word length
        if (length(arr[i]) < sw_leng) {
          sw_leng = length(arr[i])
        }
      }
      for (i=1; i<=sw_leng; i++) { # find longest common prefix
        hits = 0
        for (j=1; j<n; j++) {
          if (substr(arr[j],i,1) == substr(arr[j+1],i,1)) {
            hits++
          }
        }
        if (hits == 0) {
          break
        }
        if (hits + 1 == n) {
          lcp_leng++
        }
      }
      return(substr(str,1,lcp_leng))
    }
'
}

getGodepsForVersion() {
  local revision="$1"
  local url
  local url_fmt='https://raw.githubusercontent.com/kubernetes/client-go/%s/Godeps/Godeps.json'

  # shellcheck disable=SC2059
  url="$( printf "$url_fmt" "$revision" )"
  echo "$url"
}

renderGopkgToml() {
  local version="$1"

  curl --fail --silent "$( getGodepsForVersion "$version" )" \
    | printDepsByRev \
    | combineByRev

  # shellcheck disable=SC2059
  printf "$GOPKG_TMPL" 'k8s.io/client-go' "$version"
}

usage() {
  local me
  me="$( basename "$0" )"

  echo "${me} <version> [gopkg-toml]"
  echo ""
  echo "  Example:"
  echo "    ${me} release-5.0"
  echo "    ${me} release-4.0 ../Gopkg.toml"
}

main() {
  [ "$#" -lt 1 ] || [ "$#" -gt 2 ] && {
    usage >&2
    exit 1
  }

  local version="$1"
  local gopkgToml="${2:-"./Gopkg.toml"}"

  renderGopkgToml "$version" > "$gopkgToml"
  dep ensure
}

main "$@"
