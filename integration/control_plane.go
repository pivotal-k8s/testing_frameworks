package integration

import (
	"fmt"
	"net/url"
)

// ControlPlane is a struct that knows how to start your test control plane.
//
// Right now, that means Etcd and your APIServer. This is likely to increase in
// future.
type ControlPlane struct {
	APIServer            *APIServer
	Etcd                 *Etcd
	AdditionalComponents []ControlPlaneComponent
}

type ControlPlaneComponent interface {
	Start() error
	Stop() error
	RegisterTo(*ControlPlane)
}

// Start will start your control plane processes. To stop them, call Stop().
func (f *ControlPlane) Start() error {
	if f.Etcd == nil {
		f.Etcd = &Etcd{}
	}
	if err := f.Etcd.Start(); err != nil {
		return err
	}

	if f.APIServer == nil {
		f.APIServer = &APIServer{}
	}
	f.APIServer.EtcdURL = f.Etcd.URL
	if err := f.APIServer.Start(); err != nil {
		return err
	}

	for _, c := range f.AdditionalComponents {
		c.RegisterTo(f)
		if err := c.Start(); err != nil {
			return err
		}
	}

	return nil
}

// Stop will stop your control plane processes, and clean up their data.
func (f *ControlPlane) Stop() error {
	for _, c := range f.AdditionalComponents {
		if err := c.Stop(); err != nil {
			return nil
		}
	}

	if f.APIServer != nil {
		if err := f.APIServer.Stop(); err != nil {
			return err
		}
	}
	if f.Etcd != nil {
		if err := f.Etcd.Stop(); err != nil {
			return err
		}
	}
	return nil
}

// APIURL returns the URL you should connect to to talk to your API.
func (f *ControlPlane) APIURL() *url.URL {
	return f.APIServer.URL
}

// KubeCtl returns a pre-configured KubeCtl, ready to connect to this
// ControlPlane.
func (f *ControlPlane) KubeCtl() *KubeCtl {
	k := &KubeCtl{}
	k.Opts = append(k.Opts, fmt.Sprintf("--server=%s", f.APIURL()))
	return k
}
