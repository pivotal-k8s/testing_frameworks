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

import "net/url"

var EtcdDefaultArgs = []string{
	"--listen-peer-urls=http://localhost:0",
	"--advertise-client-urls={{ if .URL }}{{ .URL.String }}{{ end }}",
	"--listen-client-urls={{ if .URL }}{{ .URL.String }}{{ end }}",
	"--data-dir={{ .DataDir }}",
}

func DoEtcdArgDefaulting(args []string) []string {
	if len(args) != 0 {
		return args
	}

	return EtcdDefaultArgs
}

func isSecureScheme(scheme string) bool {
	// https://github.com/coreos/etcd/blob/d9deeff49a080a88c982d328ad9d33f26d1ad7b6/pkg/transport/listener.go#L53
	if scheme == "https" || scheme == "unixs" {
		return true
	}
	return false
}

func GetEtcdStartMessage(listenUrl url.URL) string {
	if isSecureScheme(listenUrl.Scheme) {
		// https://github.com/coreos/etcd/blob/a7f1fbe00ec216fcb3a1919397a103b41dca8413/embed/serve.go#L167
		return "serving client requests on " + listenUrl.Hostname()
	}

	// https://github.com/coreos/etcd/blob/a7f1fbe00ec216fcb3a1919397a103b41dca8413/embed/serve.go#L124
	return "serving insecure client requests on " + listenUrl.Hostname()
}
