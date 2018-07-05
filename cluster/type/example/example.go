// Package example holds the types needed for an example of an alternate test
// cluster implementation
package example

// NodeType is an example of additional shape description which might be needed
// by this exemplar test cluster implementation.
type NodeType struct {
	MachineType string
	AZs         []string
}
