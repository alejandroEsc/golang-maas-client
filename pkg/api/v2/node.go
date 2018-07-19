// Copyright 2016 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE File for details.

package v2

// Nodes are now known as nodes...? should reconsider this struct.
// Device represents some form of Node in maas.
type Node struct {
	// TODO: add domain
	ResourceURI string      `json:"resource_uri,omitempty"`
	SystemID    string      `json:"system_id,omitempty"`
	Hostname    string      `json:"Hostname,omitempty"`
	FQDN        string      `json:"FQDN,omitempty"`
	// Parent returns the SystemID of the Parent. Most often this will be a
	// MachineInterface.
	Parent string `json:"Parent,omitempty"`
	// Owner is the username of the user that created the Node.
	Owner       string   `json:"Owner,omitempty"`
	IPAddresses []string `json:"ip_addresses,omitempty"`
	// InterfaceSet returns all the interfaces for the NodeInterface.
	InterfaceSet []*NetworkInterface `json:"interface_set,omitempty"`
	Zone         *Zone               `json:"Zone,omitempty"`
	Tags         []string            `json:"tag_names,omitempty"`
	Type         string              `json:"node_type_name,omitempty"`
}

