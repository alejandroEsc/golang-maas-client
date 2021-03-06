// Copyright 2016 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE File for details.

package v2

import (
	"encoding/json"

	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadStaticRoutesBadSchema(t *testing.T) {
	var s StaticRoute
	err = json.Unmarshal([]byte("wat?"), &s)
	assert.Error(t, err)
}

func TestReadStaticRoutes(t *testing.T) {
	var staticRoutes []StaticRoute
	err = json.Unmarshal([]byte(staticRoutesResponse), &staticRoutes)
	assert.Nil(t, err)
	assert.Len(t, staticRoutes, 1)

	sr := staticRoutes[0]
	assert.Equal(t, sr.ID, 2)
	assert.Equal(t, sr.Metric, int(0))
	assert.Equal(t, sr.GatewayIP, "192.168.0.1")
	source := sr.Source
	assert.NotNil(t, source)
	assert.Equal(t, source.Name, "192.168.0.0/24")
	assert.Equal(t, source.CIDR, "192.168.0.0/24")
	destination := sr.Destination
	assert.NotNil(t, destination)
	assert.Equal(t, destination.Name, "Local-192")
	assert.Equal(t, destination.CIDR, "192.168.0.0/16")
}

const staticRoutesResponse = `
[
    {
        "Destination": {
            "active_discovery": false,
            "ID": 3,
            "resource_uri": "/MAAS/api/2.0/Subnets/3/",
            "allow_proxy": true,
            "rdns_mode": 2,
            "dns_servers": [
                "8.8.8.8"
            ],
            "Name": "Local-192",
            "cidr": "192.168.0.0/16",
            "Space": "Space-0",
            "VLAN": {
                "Fabric": "Fabric-1",
                "ID": 5002,
                "dhcp_on": false,
                "primary_rack": null,
                "resource_uri": "/MAAS/api/2.0/VLANs/5002/",
                "MTU": 1500,
                "fabric_id": 1,
                "secondary_rack": null,
                "Name": "untagged",
                "external_dhcp": null,
                "VID": 0
            },
            "gateway_ip": "192.168.0.1"
        },
        "Source": {
            "active_discovery": false,
            "ID": 1,
            "resource_uri": "/MAAS/api/2.0/Subnets/1/",
            "allow_proxy": true,
            "rdns_mode": 2,
            "dns_servers": [],
            "Name": "192.168.0.0/24",
            "cidr": "192.168.0.0/24",
            "Space": "Space-0",
            "VLAN": {
                "Fabric": "Fabric-0",
                "ID": 5001,
                "dhcp_on": false,
                "primary_rack": null,
                "resource_uri": "/MAAS/api/2.0/VLANs/5001/",
                "MTU": 1500,
                "fabric_id": 0,
                "secondary_rack": null,
                "Name": "untagged",
                "external_dhcp": "192.168.0.1",
                "VID": 0
            },
            "gateway_ip": null
        },
        "ID": 2,
        "resource_uri": "/MAAS/api/2.0/static-routes/2/",
        "Metric": 0,
        "gateway_ip": "192.168.0.1"
    }
]
`
