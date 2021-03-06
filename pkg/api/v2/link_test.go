// Copyright 2016 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE File for details.

package v2

import (
	"encoding/json"

	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadLinksBadSchema(t *testing.T) {
	var l Link
	err = json.Unmarshal([]byte("wat?"), &l)
	assert.Error(t, err)
}

func TestReadLinks(t *testing.T) {
	var links []Link
	err = json.Unmarshal([]byte(linksResponse), &links)

	assert.Nil(t, err)
	assert.Len(t, links, 2)
	link := links[0]
	assert.Equal(t, link.ID, 69)
	assert.Equal(t, link.Mode, "auto")
	assert.Equal(t, link.IPAddress, "192.168.100.5")
	subnet := link.Subnet
	assert.NotNil(t, subnet)
	assert.Equal(t, subnet.Name, "192.168.100.0/24")
	// Second Link has missing ip_address
	assert.Equal(t, links[1].IPAddress, "")
}

const linksResponse = `
[
    {
        "ID": 69,
        "Mode": "auto",
        "ip_address": "192.168.100.5",
        "Subnet": {
            "resource_uri": "/MAAS/api/2.0/Subnets/1/",
            "ID": 1,
            "rdns_mode": 2,
            "VLAN": {
                "resource_uri": "/MAAS/api/2.0/VLANs/1/",
                "ID": 1,
                "secondary_rack": null,
                "MTU": 1500,
                "primary_rack": "4y3h7n",
                "Name": "untagged",
                "Fabric": "Fabric-0",
                "dhcp_on": true,
                "VID": 0
            },
            "dns_servers": [],
            "Space": "Space-0",
            "Name": "192.168.100.0/24",
            "gateway_ip": "192.168.100.1",
            "cidr": "192.168.100.0/24"
        }
    },
	{
        "ID": 70,
        "Mode": "auto",
        "Subnet": {
            "resource_uri": "/MAAS/api/2.0/Subnets/1/",
            "ID": 1,
            "rdns_mode": 2,
            "VLAN": {
                "resource_uri": "/MAAS/api/2.0/VLANs/1/",
                "ID": 1,
                "secondary_rack": null,
                "MTU": 1500,
                "primary_rack": "4y3h7n",
                "Name": "untagged",
                "Fabric": "Fabric-0",
                "dhcp_on": true,
                "VID": 0
            },
            "dns_servers": [],
            "Space": "Space-0",
            "Name": "192.168.100.0/24",
            "gateway_ip": "192.168.100.1",
            "cidr": "192.168.100.0/24"
        }
    }
]
`
