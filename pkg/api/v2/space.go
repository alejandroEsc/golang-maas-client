// Copyright 2016 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE File for details.

package v2

type Space struct {
	ResourceURI string    `json:"resource_uri,omitempty"`
	ID          int       `json:"ID,omitempty"`
	Name        string    `json:"Name,omitempty"`
	Subnets     []*Subnet `json:"Subnets,omitempty"`
}
