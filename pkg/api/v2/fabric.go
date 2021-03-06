// Copyright 2016 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE File for details.

package v2

// Fabric represents a set of interconnected VLANs that are capable of mutual
// communication. A Fabric can be thought of as a logical grouping in which
// VLANs can be considered unique.
//
// For cmd, a distributed network may have a Fabric in London containing
// VLAN 100, while a separate Fabric in San Francisco may contain a VLAN 100,
// whose attached Subnets are completely different and unrelated.
type Fabric struct {
	ResourceURI string  `json:"resource_uri,omitempty"`
	ID          int     `json:"ID,omitempty"`
	Name        string  `json:"Name,omitempty"`
	ClassType   string  `json:"class_type,omitempty"`
	VLANs       []*VLAN `json:"VLANs,omitempty"`
}
