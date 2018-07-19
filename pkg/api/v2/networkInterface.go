// Copyright 2016 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE File for details.

package v2

// NetworkInterface represents a physical or virtual network interface on a MachineInterface.
type NetworkInterface struct {
	ResourceURI  string   `json:"resource_uri,omitempty"`
	ID           int      `json:"ID,omitempty"`
	Name         string   `json:"Name,omitempty"`
	Type         string   `json:"type,omitempty"`
	Enabled      bool     `json:"Enabled,omitempty"`
	Tags         []string `json:"Tags,omitempty"`
	VLAN         *VLAN    `json:"VLAN,omitempty"`
	Links        []*Link  `json:"Links,omitempty"`
	MACAddress   string   `json:"mac_address,omitempty"`
	EffectiveMTU int      `json:"effective_mtu,omitempty"`
	Parents      []string `json:"Parents,omitempty"`
	Children     []string `json:"Children,omitempty"`
}

func (i *NetworkInterface) updateFrom(other *NetworkInterface) {
	i.ResourceURI = other.ResourceURI
	i.ID = other.ID
	i.Name = other.Name
	i.Type = other.Type
	i.Enabled = other.Enabled
	i.Tags = other.Tags
	i.VLAN = other.VLAN
	i.Links = other.Links
	i.MACAddress = other.MACAddress
	i.EffectiveMTU = other.EffectiveMTU
	i.Parents = other.Parents
	i.Children = other.Children
}

func (i *NetworkInterface) linkForSubnet(st *Subnet) *Link {
	for _, link := range i.Links {
		if s := link.Subnet; s != nil && s.ID == st.ID {
			return link
		}
	}
	return nil
}
