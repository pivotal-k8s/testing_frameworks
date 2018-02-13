package internal

var VirtualKubeletDefaultArgs = []string{
	"--provider=mock",
	"--kubeconfig={{ .ConfDir }}/kube.conf",
}

func DoVirtualKubeletArgDefaulting(args []string) []string {
	if len(args) != 0 {
		return args
	}

	return VirtualKubeletDefaultArgs
}
