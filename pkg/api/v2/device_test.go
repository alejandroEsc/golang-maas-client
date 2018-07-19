// Copyright 2016 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE File for details.

package v2

import (
	"net/http"

	"encoding/json"

	"testing"

	"github.com/alejandroEsc/golang-maas-client/pkg/api/util"
	"github.com/stretchr/testify/assert"
)

func TestReadDeviceBadSchema(t *testing.T) {
	var d Device
	err = json.Unmarshal([]byte("wat?"), &d)
	assert.Error(t, err)
}

func TestReadDevices(t *testing.T) {
	var devices []Device
	err = json.Unmarshal([]byte(devicesResponse), &devices)
	assert.Nil(t, err)

	assert.Len(t, devices, 1)

	device := devices[0]
	assert.Equal(t, "4y3haf", device.SystemID)
	assert.Equal(t, "furnacelike-brittney", device.Hostname)
	assert.Equal(t, "furnacelike-brittney.maas", device.FQDN)
	assert.EqualValues(t, []string{"192.168.100.11"}, device.IPAddresses)
	zone := device.Zone
	assert.NotNil(t, zone)
	assert.Equal(t, "default", zone.Name)
}

func TestDeviceDelete(t *testing.T) {
	server, controller := createTestServerController(t)
	// Successful delete is 204 - StatusNoContent
	server.AddGetResponse("/api/2.0/devices/", http.StatusOK, devicesResponse)
	devices, err := controller.Devices(DevicesArgs{})
	assert.Nil(t, err)
	assert.Len(t, devices, 1)

	server.AddDeleteResponse(devices[0].ResourceURI, http.StatusNoContent, "")
	err = controller.DeleteDevice(devices[0])
	assert.Nil(t, err)
}

func TestDeviceDelete404(t *testing.T) {
	server, controller := createTestServerController(t)
	server.AddGetResponse("/api/2.0/devices/", http.StatusOK, devicesResponse)
	devices, err := controller.Devices(DevicesArgs{})
	assert.Nil(t, err)
	assert.Len(t, devices, 1)
	// No Path, so 404
	err = controller.DeleteDevice(devices[0])
	assert.True(t, util.IsNoMatchError(err))
}

func TestDeviceDeleteForbidden(t *testing.T) {
	server, controller := createTestServerController(t)
	server.AddGetResponse("/api/2.0/devices/", http.StatusOK, devicesResponse)
	devices, err := controller.Devices(DevicesArgs{})
	assert.Nil(t, err)
	assert.Len(t, devices, 1)

	server.AddDeleteResponse(devices[0].ResourceURI, http.StatusForbidden, "")
	err = controller.DeleteDevice(devices[0])
	assert.True(t, util.IsPermissionError(err))
}

func TestDeviceDeleteUnknown(t *testing.T) {
	server, controller := createTestServerController(t)
	server.AddGetResponse("/api/2.0/devices/", http.StatusOK, devicesResponse)
	devices, err := controller.Devices(DevicesArgs{})
	assert.Nil(t, err)
	assert.Len(t, devices, 1)

	server.AddDeleteResponse(devices[0].ResourceURI, http.StatusConflict, "")
	assert.Nil(t, err)
	assert.Len(t, devices, 1)
	err = controller.DeleteDevice(devices[0])
	assert.True(t, util.IsUnexpectedError(err))
}

const (
	deviceResponse = `
    {
        "Zone": {
            "Description": "",
            "resource_uri": "/maas/api/2.0/zones/default/",
            "Name": "default"
        },
        "domain": {
            "resource_record_count": 0,
            "resource_uri": "/maas/api/2.0/domains/0/",
            "authoritative": true,
            "Name": "maas",
            "ttl": null,
            "ID": 0
        },
        "Device_type_name": "DeviceInterface",
        "address_ttl": null,
        "Hostname": "furnacelike-brittney",
        "Device_type": 1,
        "resource_uri": "/maas/api/2.0/Devices/4y3haf/",
        "ip_addresses": ["192.168.100.11"],
        "Owner": "thumper",
        "tag_names": [],
        "FQDN": "furnacelike-brittney.maas",
        "system_id": "4y3haf",
        "Parent": "4y3ha3",
        "interface_set": [
            {
                "resource_uri": "/maas/api/2.0/Devices/4y3haf/interfaces/48/",
                "type": "physical",
                "mac_address": "78:f0:f1:16:a7:46",
                "params": "",
                "discovered": null,
                "effective_mtu": 1500,
                "ID": 48,
                "Children": [],
                "Links": [],
                "Name": "eth0",
                "VLAN": {
                    "secondary_rack": null,
                    "dhcp_on": true,
                    "Fabric": "Fabric-0",
                    "MTU": 1500,
                    "primary_rack": "4y3h7n",
                    "resource_uri": "/maas/api/2.0/VLANs/1/",
                    "external_dhcp": null,
                    "Name": "untagged",
                    "ID": 1,
                    "VID": 0
                },
                "Tags": [],
                "Parents": [],
                "Enabled": true
            },
            {
                "resource_uri": "/maas/api/2.0/Devices/4y3haf/interfaces/49/",
                "type": "physical",
                "mac_address": "15:34:d3:2d:f7:a7",
                "params": {},
                "discovered": null,
                "effective_mtu": 1500,
                "ID": 49,
                "Children": [],
                "Links": [
                    {
                        "Mode": "link_up",
                        "ID": 101
                    }
                ],
                "Name": "eth1",
                "VLAN": {
                    "secondary_rack": null,
                    "dhcp_on": true,
                    "Fabric": "Fabric-0",
                    "MTU": 1500,
                    "primary_rack": "4y3h7n",
                    "resource_uri": "/maas/api/2.0/VLANs/1/",
                    "external_dhcp": null,
                    "Name": "untagged",
                    "ID": 1,
                    "VID": 0
                },
                "Tags": [],
                "Parents": [],
                "Enabled": true
            }
        ]
    }
    `
	devicesResponse = "[" + deviceResponse + "]"
)
