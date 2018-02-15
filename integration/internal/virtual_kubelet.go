package internal

var VirtualKubeletDefaultArgs = []string{
	"--provider=mock",
	"--kubeconfig={{ .Dir }}/kube.conf",
}

func DoVirtualKubeletArgDefaulting(args []string) []string {
	if len(args) != 0 {
		return args
	}

	return VirtualKubeletDefaultArgs
}
