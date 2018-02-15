#!/usr/bin/env bash
set -eu

# Use DEBUG=1 ./scripts/download-binaries.sh to get debug output
quiet="-s"
[[ -z "${DEBUG:-""}" ]] || {
  set -x
  quiet=""
}

logEnd() {
  local msg='done.'
  [ "$1" -eq 0 ] || msg='Error downloading assets'
  echo "$msg"
}
trap 'logEnd $?' EXIT

# Use BASE_URL=https://my/binaries/url ./scripts/download-binaries to download
# from a different bucket
: "${BASE_URL:="https://storage.googleapis.com/k8s-c10s-test-binaries"}"

test_framework_dir="$(cd "$(dirname "$0")/.." ; pwd)"
os="$(uname -s)"
os_lowercase="$(echo "$os" | tr '[:upper:]' '[:lower:]' )"
arch="$(uname -m)"

dest_dir="${1:-"${test_framework_dir}/assets/bin"}"
etcd_dest="${dest_dir}/etcd"
kube_apiserver_dest="${dest_dir}/kube-apiserver"
kube_controller_manager_dest="${dest_dir}/kube-controller-manager"
kube_scheduler_dest="${dest_dir}/kube-scheduler"
virtual_kubelet_dest="${dest_dir}/virtual-kubelet"
kubectl_dest="${dest_dir}/kubectl"

echo "About to download a couple of binaries. This might take a while..."

curl $quiet "${BASE_URL}/etcd-${os}-${arch}" --output "$etcd_dest"
curl $quiet "${BASE_URL}/kube-apiserver-${os}-${arch}" --output "$kube_apiserver_dest"
curl $quiet "${BASE_URL}/kube-controller-manager-${os}-${arch}" --output "$kube_controller_manager_dest"
curl $quiet "${BASE_URL}/kube-scheduler-${os}-${arch}" --output "$kube_scheduler_dest"
curl $quiet "${BASE_URL}/virtual-kubelet-${os}-${arch}" --output "$virtual_kubelet_dest"

kubectl_version="$(curl $quiet https://storage.googleapis.com/kubernetes-release/release/stable.txt)"
kubectl_url="https://storage.googleapis.com/kubernetes-release/release/${kubectl_version}/bin/${os_lowercase}/amd64/kubectl"
curl $quiet "$kubectl_url" --output "$kubectl_dest"

chmod +x \
  "$etcd_dest" \
  "$kube_apiserver_dest" \
  "$kube_controller_manager_dest" \
  "$kube_scheduler_dest" \
  "$virtual_kubelet_dest" \
  "$kubectl_dest"

echo    "#   ${dest_dir}"
echo    "# versions:"
echo -n "#   etcd:                    "; "$etcd_dest" --version | head -n 1
echo -n "#   kube-apiserver:          "; "$kube_apiserver_dest" --version
echo -n "#   kube-controller-manager: "; "$kube_controller_manager_dest" --version
echo -n "#   kube-scheduler:          "; "$kube_scheduler_dest" --version
echo -n "#   virtual-kubelet:         "; "$virtual_kubelet_dest" --provider=doesNotMatter version
echo -n "#   kubectl:                 "; "$kubectl_dest" version --client --short
