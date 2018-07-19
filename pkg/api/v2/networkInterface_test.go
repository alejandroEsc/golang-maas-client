// Copyright 2016 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE File for details.

package v2

import (
	"net/http"

	"encoding/json"

	"testing"

	"github.com/alejandroEsc/golang-maas-client/pkg/api/client"
	"github.com/alejandroEsc/golang-maas-client/pkg/api/util"
	"github.com/juju/errors"
	"github.com/stretchr/testify/assert"
)

func TestReadNetworkInterfacesBadSchema(t *testing.T) {
	var b NetworkInterface
	err = json.Unmarshal([]byte("wat?"), &b)
	assert.Error(t, err)
}

func TestReadNetworkInterface(t *testing.T) {
	var iface NetworkInterface
	err = json.Unmarshal([]byte(interfaceResponse), &iface)
	assert.Nil(t, err)
	checkNetworkInterface(t, &iface)
}

func TestReadNetworkInterfacesNulls(t *testing.T) {
	var iface NetworkInterface
	err = json.Unmarshal([]byte(interfaceNullsResponse), &iface)

	assert.Nil(t, err)

	assert.Equal(t, "", iface.MACAddress)
	assert.Equal(t, []string(nil), iface.Tags)
	assert.Nil(t, iface.VLAN)
}

func TestReadNeworkInterfaces(t *testing.T) {
	var iface []NetworkInterface
	err = json.Unmarshal([]byte(interfacesResponse), &iface)
	assert.Nil(t, err)
	assert.Len(t, iface, 1)
	checkNetworkInterface(t, &iface[0])
}

func TestNetworkInterfaceLinkSubnetValidates(t *testing.T) {
	server, iface, controller := getServerNewInterfaceAndController(t)
	defer server.Close()

	err := controller.LinkSubnet(iface, LinkSubnetArgs{})
	assert.True(t, errors.IsNotValid(err))
	assert.Equal(t, err.Error(), "missing Mode not valid")
}

func TestNetworkInterfaceLinkSubnetGood(t *testing.T) {
	server, iface, controller := getServerNewInterfaceAndController(t)
	// The changed information is there just for the test to show that the response
	// is parsed and the interface updated
	response := util.UpdateJSONMap(t, interfaceResponse, map[string]interface{}{
		"Name": "eth42",
	})
	defer server.Close()

	server.AddPostResponse(iface.ResourceURI+"?op=link_subnet", http.StatusOK, response)
	args := LinkSubnetArgs{
		Mode:           LinkModeStatic,
		Subnet:         &Subnet{ID: 42},
		IPAddress:      "10.10.10.10",
		DefaultGateway: true,
	}
	err := controller.LinkSubnet(iface, args)
	assert.Nil(t, err)
	assert.Equal(t, iface.Name, "eth42")

	request := server.LastRequest()
	form := request.PostForm
	assert.Equal(t, form.Get("Mode"), "STATIC")
	assert.Equal(t, form.Get("Subnet"), "42")
	assert.Equal(t, form.Get("ip_address"), "10.10.10.10")
	assert.Equal(t, form.Get("default_gateway"), "true")
}

func TestNetworkInterfaceLinkSubnetMissing(t *testing.T) {
	server, iface, controller := getServerNewInterfaceAndController(t)
	defer server.Close()

	args := LinkSubnetArgs{
		Mode:   LinkModeStatic,
		Subnet: &Subnet{ID: 42},
	}
	err := controller.LinkSubnet(iface, args)
	assert.True(t, util.IsBadRequestError(err))
}

func TestNetworkInterfaceLinkSubnetForbidden(t *testing.T) {
	server, iface, controller := getServerNewInterfaceAndController(t)
	server.AddPostResponse(iface.ResourceURI+"?op=link_subnet", http.StatusForbidden, "bad user")
	defer server.Close()

	args := LinkSubnetArgs{
		Mode:   LinkModeStatic,
		Subnet: &Subnet{ID: 42},
	}
	err := controller.LinkSubnet(iface, args)
	assert.True(t, util.IsPermissionError(err))
	assert.Equal(t, err.Error(), "bad user")
}

func TestNetworkInterfaceLinkSubnetNoAddressesAvailable(t *testing.T) {
	server, iface, controller := getServerNewInterfaceAndController(t)
	server.AddPostResponse(iface.ResourceURI+"?op=link_subnet", http.StatusServiceUnavailable, "no addresses")
	defer server.Close()

	args := LinkSubnetArgs{
		Mode:   LinkModeStatic,
		Subnet: &Subnet{ID: 42},
	}
	err := controller.LinkSubnet(iface, args)
	assert.True(t, util.IsCannotCompleteError(err))
	assert.Equal(t, err.Error(), "no addresses")
}

func TestNetworkInterfaceLinkSubnetUnknown(t *testing.T) {
	server, iface, controller := getServerNewInterfaceAndController(t)
	server.AddPostResponse(iface.ResourceURI+"?op=link_subnet", http.StatusMethodNotAllowed, "wat?")
	defer server.Close()

	args := LinkSubnetArgs{
		Mode:   LinkModeStatic,
		Subnet: &Subnet{ID: 42},
	}
	err := controller.LinkSubnet(iface, args)
	assert.True(t, util.IsUnexpectedError(err))
	assert.Equal(t, err.Error(), "unexpected: ServerError: 405 Method Not Allowed (wat?)")
}

func TestNetworkInterfaceUnlinkSubnetValidates(t *testing.T) {
	server, iface, controller := getServerNewInterfaceAndController(t)
	defer server.Close()

	err := controller.UnlinkSubnet(iface, nil)
	assert.True(t, errors.IsNotValid(err))
	assert.Equal(t, err.Error(), "missing Subnet not valid")
}

func TestNetworkInterfaceUnlinkSubnetNotLinked(t *testing.T) {
	server, iface, controller := getServerNewInterfaceAndController(t)
	defer server.Close()

	err := controller.UnlinkSubnet(iface, &Subnet{ID: 42})
	assert.True(t, errors.IsNotValid(err))
	assert.Equal(t, err.Error(), "unlinked Subnet not valid")
}

func TestNetworkInterfaceUnlinkSubnetGood(t *testing.T) {
	server, iface, controller := getServerNewInterfaceAndController(t)
	// The changed information is there just for the test to show that the response
	// is parsed and the interface updated
	response := util.UpdateJSONMap(t, interfaceResponse, map[string]interface{}{
		"Name": "eth42",
	})
	server.AddPostResponse(iface.ResourceURI+"?op=unlink_subnet", http.StatusOK, response)
	defer server.Close()

	err := controller.UnlinkSubnet(iface, &Subnet{ID: 1})
	assert.Nil(t, err)
	assert.Equal(t, iface.Name, "eth42")

	request := server.LastRequest()
	form := request.PostForm
	// The Link ID that contains Subnet 1 has an internal ID of 69.
	assert.Equal(t, form.Get("ID"), "69")
}

func TestNetworkInterfaceUnlinkSubnetMissing(t *testing.T) {
	server, iface, controller := getServerNewInterfaceAndController(t)
	defer server.Close()

	err := controller.UnlinkSubnet(iface, &Subnet{ID: 1})
	assert.True(t, util.IsBadRequestError(err))
}

func TestNetworkInterfaceUnlinkSubnetForbidden(t *testing.T) {
	server, iface, controller := getServerNewInterfaceAndController(t)
	server.AddPostResponse(iface.ResourceURI+"?op=unlink_subnet", http.StatusForbidden, "bad user")
	defer server.Close()

	err := controller.UnlinkSubnet(iface, &Subnet{ID: 1})
	assert.True(t, util.IsPermissionError(err))
	assert.Equal(t, err.Error(), "bad user")
}

func TestNetworkInterfaceUnlinkSubnetUnknown(t *testing.T) {
	server, iface, controller := getServerNewInterfaceAndController(t)
	server.AddPostResponse(iface.ResourceURI+"?op=unlink_subnet", http.StatusMethodNotAllowed, "wat?")
	defer server.Close()

	err := controller.UnlinkSubnet(iface, &Subnet{ID: 1})
	assert.True(t, util.IsUnexpectedError(err))
	assert.Equal(t, err.Error(), "unexpected: ServerError: 405 Method Not Allowed (wat?)")
}

func TestNetworkInterfaceUpdateNoChangeNoRequest(t *testing.T) {
	server, iface, controller := getServerNewInterfaceAndController(t)
	defer server.Close()

	count := server.RequestCount()
	err := controller.UpdateNetworkInterface(iface, UpdateInterfaceArgs{})
	assert.Nil(t, err)
	assert.Equal(t, server.RequestCount(), count)
}

func TestNetworkInterfaceUpdateMissing(t *testing.T) {
	_, iface, controller := getServerNewInterfaceAndController(t)
	err := controller.UpdateNetworkInterface(iface, UpdateInterfaceArgs{Name: "eth2"})
	assert.True(t, util.IsNoMatchError(err))
}

func TestInterfaceUpdateForbidden(t *testing.T) {
	server, iface, controller := getServerNewInterfaceAndController(t)
	server.AddPutResponse(iface.ResourceURI, http.StatusForbidden, "bad user")
	defer server.Close()

	err := controller.UpdateNetworkInterface(iface, UpdateInterfaceArgs{Name: "eth2"})
	assert.True(t, util.IsPermissionError(err))
	assert.Equal(t, err.Error(), "bad user")
}

func TestNetworkInterfaceUpdateUnknown(t *testing.T) {
	server, iface, controller := getServerNewInterfaceAndController(t)
	server.AddPutResponse(iface.ResourceURI, http.StatusMethodNotAllowed, "wat?")
	defer server.Close()

	err := controller.UpdateNetworkInterface(iface, UpdateInterfaceArgs{Name: "eth2"})
	assert.True(t, util.IsUnexpectedError(err))
	assert.Equal(t, err.Error(), "unexpected: ServerError: 405 Method Not Allowed (wat?)")
}

func TestNetworkInterfaceUpdateGood(t *testing.T) {
	server, iface, controller := getServerNewInterfaceAndController(t)
	// The changed information is there just for the test to show that the response
	// is parsed and the interface updated
	defer server.Close()

	response := util.UpdateJSONMap(t, interfaceResponse, map[string]interface{}{
		"Name": "eth42",
	})
	server.AddPutResponse(iface.ResourceURI, http.StatusOK, response)
	args := UpdateInterfaceArgs{
		Name:       "eth42",
		MACAddress: "c3-52-51-b4-50-cd",
		VLAN:       VLAN{ID: 13},
	}
	err := controller.UpdateNetworkInterface(iface, args)
	assert.Nil(t, err)
	assert.Equal(t, iface.Name, "eth42")

	request := server.LastRequest()
	form := request.PostForm
	assert.Equal(t, form.Get("Name"), "eth42")
	assert.Equal(t, form.Get("mac_address"), "c3-52-51-b4-50-cd")
	assert.Equal(t, form.Get("VLAN"), "13")
}

func getServerNewInterfaceAndController(t *testing.T) (*client.SimpleTestServer, *NetworkInterface, *Controller) {
	server, controller := createTestServerController(t)
	server.AddGetResponse("/api/2.0/nodes/", http.StatusOK, nodesResponse)

	nodes, err := controller.Nodes(NodesArgs{})
	assert.Nil(t, err)
	node := nodes[0]
	server.AddPostResponse(node.ResourceURI+"interfaces/?op=create_physical", http.StatusOK, interfaceResponse)
	iface, err := controller.CreateInterface(&node, minimalCreateInterfaceArgs())
	assert.Nil(t, err)
	return server, iface, controller
}

func checkNetworkInterface(t *testing.T, iface *NetworkInterface) {
	assert.Equal(t, 40, iface.ID)
	assert.Equal(t, "eth0", iface.Name)
	assert.Equal(t, iface.Type, "physical")
	assert.True(t, iface.Enabled)
	assert.Equal(t, iface.Tags, []string{"foo", "bar"})

	assert.Equal(t, iface.MACAddress, "52:54:00:c9:6a:45")
	assert.Equal(t, iface.EffectiveMTU, 1500)

	assert.Equal(t, iface.Parents, []string{"bond0"})
	assert.Equal(t, iface.Children, []string{"eth0.1", "eth0.2"})

	vlan := iface.VLAN
	assert.NotNil(t, vlan)
	assert.NotNil(t, vlan.Name)

	assert.Equal(t, vlan.Name, "untagged")

	links := iface.Links
	assert.Len(t, links, 1)
	assert.Equal(t, links[0].ID, 69)
}

const (
	interfacesResponse = "[" + interfaceResponse + "]"
	interfaceResponse  = `
{
    "effective_mtu": 1500,
    "mac_address": "52:54:00:c9:6a:45",
    "Children": ["eth0.1", "eth0.2"],
    "discovered": [],
    "params": "some params",
    "VLAN": {
        "resource_uri": "/maas/api/2.0/VLANs/1/",
        "ID": 1,
        "secondary_rack": null,
        "MTU": 1500,
        "primary_rack": "4y3h7n",
        "Name": "untagged",
        "Fabric": "Fabric-0",
        "dhcp_on": true,
        "VID": 0
    },
    "Name": "eth0",
    "Enabled": true,
    "Parents": ["bond0"],
    "ID": 40,
    "type": "physical",
    "resource_uri": "/maas/api/2.0/nodes/4y3ha6/interfaces/40/",
    "Tags": ["foo", "bar"],
    "Links": [
        {
            "ID": 69,
            "Mode": "auto",
            "Subnet": {
                "resource_uri": "/maas/api/2.0/Subnets/1/",
                "ID": 1,
                "rdns_mode": 2,
                "VLAN": {
                    "resource_uri": "/maas/api/2.0/VLANs/1/",
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
}
`
	interfaceNullsResponse = `
{
    "effective_mtu": 1500,
    "mac_address": null,
    "Children": ["eth0.1", "eth0.2"],
    "discovered": [],
    "params": "some params",
    "VLAN": null,
    "Name": "eth0",
    "Enabled": true,
    "Parents": ["bond0"],
    "ID": 40,
    "type": "physical",
    "resource_uri": "/maas/api/2.0/nodes/4y3ha6/interfaces/40/",
    "Tags": null,
    "Links": [
        {
            "ID": 69,
            "Mode": "auto",
            "Subnet": {
                "resource_uri": "/maas/api/2.0/Subnets/1/",
                "ID": 1,
                "rdns_mode": 2,
                "VLAN": {
                    "resource_uri": "/maas/api/2.0/VLANs/1/",
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
}
`
)
