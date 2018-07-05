package lightweight

import (
	"sigs.k8s.io/testing_frameworks/cluster"
	"sigs.k8s.io/testing_frameworks/cluster/type/lightweight"
)

type Node struct {
	ClusterConfig cluster.Config
	KubeletType   lightweight.KubeletType
}

func (node *Node) Start() error {
	// TODO: not implemented yet.
	return nil
}

func (node *Node) Stop() error {
	// TODO: not implemented yey.
	return nil
}
