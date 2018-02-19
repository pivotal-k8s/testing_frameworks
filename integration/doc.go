/*

Package integration implements an integration testing framework for kubernetes.

It provides components for standing up a kubernetes API, against which you can test a
kubernetes client, or other kubernetes components. The lifecycle of the components
needed to provide this API is managed by this framework.

Quickstart

If you want to test a kubernetes client against the latest kubernetes APIServer
and Etcd, you can use `./scripts/download-binaries.sh` to download APIServer and
Etcd binaries for your platform. Then add something like the
following to your tests:

	cp := &integration.ControlPlane{}
	cp.Start()
	kubeCtl, err := cp.KubeCtl()
	stdout, stderr, err := kubeCtl.Run("get", "pods")
	// You can check on err, stdout & stderr and build up
	// your tests
	cp.Stop()

Main Components

Currently the framework provides the following components:

ControlPlane: The ControlPlane wraps Etcd & APIServer (see below) and wires
them together correctly. A ControlPlane can be stopped & started and can
provide the URL to connect to the API. The ControlPlane can also be asked for a
KubeCtl which is already correctly configured for this ControlPlane. The
ControlPlane is a good entry point for default setups.

Etcd: Manages an Etcd binary, which can be started, stopped and connected to.
By default Etcd will listen on a random port for http connections and will
create a temporary directory for its data. To configure it differently, see the
Etcd type documentation below.

APIServer: Manages an Kube-APIServer binary, which can be started, stopped and
connected to. By default APIServer will listen on a random port for http
connections and will create a temporary directory to store the (auto-generated)
certificates.  To configure it differently, see the APIServer type
documentation below.

KubeCtl: Wraps around a `kubectl` binary and can `Run(...)` arbitrary commands
against a kubernetes control plane.

Additional Components

A control plane can be configured with more components then just Etcd and the
APIServer.

To do so, an additional component needs to implement the
`ControlPlaneComponent` interface and the control plane needs to be configured
to use this additional component.

This package ships with some additional components (ControllerManager,
Scheduler, VirtualKubelet). Those are not added to the control plane by
default. If you wish to use one or more of those, set up the control plane like
that:

	scheduler := &integration.Scheduler{}
	vKubelet := &integration.VirtualKubelet{}
	cp := &integration.ControlPlane{
		AdditionalComponents: []integration.ControlPlaneComponent{
			scheduler,
			vKubelet,
		}
	}
	cp.Start()
	// exercise the control plane
	cp.Stop()

Binaries

The framework uses the same mechanism to determine which binaries to
use when starting components.

For convenience this framework ships with
`${FRAMEWORK_DIR}/scripts/download-binaries.sh` which can be used to download
pre-compiled versions of the needed binaries and place them in the default
location (`${FRAMEWORK_DIR}/assets/bin/`).

1. By default the framework expects to find the binaries in the directory
`${FRAMEWORK_DIR}/assets/bin/`.

2. The path to a binary can be configured using an environment variable:

	$TEST_ASSET_KUBE_APISERVER
	$TEST_ASSET_ETCD
	$TEST_ASSET_KUBECTL
	$TEST_ASSET_KUBE_SCHEDULER
	$TEST_ASSET_KUBE_CONTROLLER_MANAGER
	$TEST_ASSET_VIRTUAL_KUBELET

3. Alternatively, the path to a binary can be configured using the `Path` field
on the component struct.
For example:

	myEtcd := &integration.Etcd{
		Path: "/some/other/etcd",
	}
	cp := &integration.ControlPlane{
		Etcd: myEtcd,
	}
	cp.Start()

Custom Arguments for Components

All components will start with default arguments. These can be customised by
setting the `Args` field on the component struct.

The default arguments are available as package variables. Therefore, you can
use them in your test set up, e.g. appending to them.

All arguments are interpreted as go templates. Those templates have access to
some fields depending on the struct. The templates are evaluated just before
the binary is executed.

	etcdArgs := append(integration.APIServerDefaultArgs, []string{
		"--debug",
		"--wal-dir={{ .Dir }}/etcd_wals/",
	})

	etcd := &Etcd{
		Args:    etcdArgs,
		DataDir: "/my/special/data/dir",
	}

*/
package integration
