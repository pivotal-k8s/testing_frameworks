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

import (
	"fmt"
	"net"
)

// AddressManager allocates a new address (interface & port) a process
// can bind and keeps track of that.
type AddressManager struct {
	port int
	host string
}

// Initialize returns a address a process can listen on. It returns
// a tuple consisting of a free port and the hostname resolved to its IP.
func (d *AddressManager) Initialize() (port int, resolvedHost string, err error) {
	if d.port != 0 {
		return 0, "", fmt.Errorf("this AddressManager is already initialized")
	}
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return
	}
	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return
	}
	d.port = l.Addr().(*net.TCPAddr).Port
	defer func() {
		err = l.Close()
	}()
	d.host = addr.IP.String()
	return d.port, d.host, nil
}

// Port returns the port that this AddressManager is managing. Port returns an
// error if this AddressManager has not yet been initialized.
func (d *AddressManager) Port() (int, error) {
	if d.port == 0 {
		return 0, fmt.Errorf("this AdressManager is not initialized yet")
	}
	return d.port, nil
}

// Host returns the host that this AddressManager is managing. Host returns an
// error if this AddressManager has not yet been initialized.
func (d *AddressManager) Host() (string, error) {
	if d.host == "" {
		return "", fmt.Errorf("this AdressManager is not initialized yet")
	}
	return d.host, nil
}
