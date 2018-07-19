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

func TestReadNodesBadSchema(t *testing.T) {
	var d Node
	err = json.Unmarshal([]byte("wat?"), &d)
	assert.Error(t, err)
}

func TestReadNodes(t *testing.T) {
	var devices []Node
	err = json.Unmarshal([]byte(nodesResponse), &devices)
	assert.Nil(t, err)

	assert.Len(t, devices, 1)

	device := devices[0]
	assert.Equal(t, device.SystemID, "4y3haf")
	assert.Equal(t, device.Hostname, "furnacelike-brittney")
	assert.Equal(t, device.FQDN, "furnacelike-brittney.maas")
	assert.EqualValues(t, device.IPAddresses, []string{"192.168.100.11"})
	zone := device.Zone
	assert.NotNil(t, zone)
	assert.Equal(t, zone.Name, "default")
}

func TestNodeInterfaceSet(t *testing.T) {
	server, node, _ := getServeNodeAndController(t)
	server.AddGetResponse(node.ResourceURI+"interfaces/", http.StatusOK, interfacesResponse)
	defer server.Close()

	ifaces := node.InterfaceSet
	assert.Len(t, ifaces, 2)
}

func TestNodeCreateInterfaceValidates(t *testing.T) {
	server, node, controller := getServeNodeAndController(t)
	defer server.Close()

	_, err := controller.CreateInterface(node, CreateNodeNetworkInterfaceArgs{})
	assert.True(t, errors.IsNotValid(err))
}

func TestNodeCreateInterface(t *testing.T) {
	server, node, controller := getServeNodeAndController(t)
	server.AddPostResponse(node.ResourceURI+"interfaces/?op=create_physical", http.StatusOK, interfaceResponse)
	defer server.Close()

	args := CreateNodeNetworkInterfaceArgs{
		Name:       "eth43",
		MACAddress: "some-mac-address",
		VLAN:       VLAN{ID: 33},
		Tags:       []string{"foo", "bar"},
	}

	iface, err := controller.CreateInterface(node, args)
	assert.Nil(t, err)
	assert.NotNil(t, iface)

	request := server.LastRequest()
	form := request.PostForm
	assert.Equal(t, form.Get("Name"), "eth43")
	assert.Equal(t, form.Get("mac_address"), "some-mac-address")
	assert.Equal(t, form.Get("VLAN"), "33")
	assert.Equal(t, form.Get("Tags"), "foo,bar")
}

func minimalCreateInterfaceArgs() CreateNodeNetworkInterfaceArgs {
	return CreateNodeNetworkInterfaceArgs{
		Name:       "eth43",
		MACAddress: "some-mac-address",
		VLAN:       VLAN{ID: 33},
	}
}

func TestNodeCreateInterfaceNotFound(t *testing.T) {
	server, node, controller := getServeNodeAndController(t)
	server.AddPostResponse(node.ResourceURI+"interfaces/?op=create_physical", http.StatusNotFound, "can't find Node")
	defer server.Close()
	_, err := controller.CreateInterface(node, minimalCreateInterfaceArgs())
	assert.True(t, util.IsBadRequestError(err))
	assert.Equal(t, err.Error(), "can't find Node")
}

func TestCreateInterfaceConflict(t *testing.T) {
	server, node, controller := getServeNodeAndController(t)
	server.AddPostResponse(node.ResourceURI+"interfaces/?op=create_physical", http.StatusConflict, "Node not allocated")
	defer server.Close()
	_, err := controller.CreateInterface(node, minimalCreateInterfaceArgs())
	assert.True(t, util.IsBadRequestError(err))
	assert.Equal(t, err.Error(), "Node not allocated")
}

func TestCreateInterfaceForbidden(t *testing.T) {
	server, node, controller := getServeNodeAndController(t)
	server.AddPostResponse(node.ResourceURI+"interfaces/?op=create_physical", http.StatusForbidden, "Node not yours")
	defer server.Close()
	_, err := controller.CreateInterface(node, minimalCreateInterfaceArgs())
	assert.True(t, util.IsPermissionError(err))
	assert.Equal(t, err.Error(), "Node not yours")
}

func TestCreateInterfaceServiceUnavailable(t *testing.T) {
	server, node, controller := getServeNodeAndController(t)
	server.AddPostResponse(node.ResourceURI+"interfaces/?op=create_physical", http.StatusServiceUnavailable, "no ip addresses available")
	defer server.Close()
	_, err := controller.CreateInterface(node, minimalCreateInterfaceArgs())
	assert.True(t, util.IsCannotCompleteError(err))
	assert.Equal(t, err.Error(), "no ip addresses available")
}

func TestNodeCreateInterfaceUnknown(t *testing.T) {
	server, node, controller := getServeNodeAndController(t)
	server.AddPostResponse(node.ResourceURI+"interfaces/?op=create_physical", http.StatusMethodNotAllowed, "wat?")
	defer server.Close()
	_, err := controller.CreateInterface(node, minimalCreateInterfaceArgs())
	assert.True(t, util.IsUnexpectedError(err))
	assert.Equal(t, err.Error(), "unexpected: ServerError: 405 Method Not Allowed (wat?)")
}

func TestNodeDelete(t *testing.T) {
	server, node, controller := getServeNodeAndController(t)
	// Successful delete is 204 - StatusNoContent
	server.AddDeleteResponse(node.ResourceURI, http.StatusNoContent, "")
	defer server.Close()
	err := controller.DeleteNode(node)
	assert.Nil(t, err)
}

func TestNodeDelete404(t *testing.T) {
	_, node, controller := getServeNodeAndController(t)
	// No Path, so 404
	defer server.Close()
	err := controller.DeleteNode(node)
	assert.True(t, util.IsNoMatchError(err))
}

func TestNodeDeleteForbidden(t *testing.T) {
	server, node, controller := getServeNodeAndController(t)
	server.AddDeleteResponse(node.ResourceURI, http.StatusForbidden, "")
	defer server.Close()
	err := controller.DeleteNode(node)
	assert.True(t, util.IsPermissionError(err))
}

func TestNodeDeleteUnknown(t *testing.T) {
	server, node, controller := getServeNodeAndController(t)
	server.AddDeleteResponse(node.ResourceURI, http.StatusConflict, "")
	defer server.Close()
	err := controller.DeleteNode(node)
	assert.True(t, util.IsUnexpectedError(err))
}

func getServeNodeAndController(t *testing.T) (*client.SimpleTestServer, *Node, *Controller) {
	server, controller := createTestServerController(t)
	server.AddGetResponse("/api/2.0/nodes/", http.StatusOK, nodesResponse)

	devices, err := controller.Nodes(NodesArgs{})
	assert.Nil(t, err)
	assert.Len(t, devices, 1)
	return server, &devices[0], controller
}

const (
	nodeResponse = `
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
        "node_type_name": "NodeInterface",
        "address_ttl": null,
        "Hostname": "furnacelike-brittney",
        "node_type": 1,
        "resource_uri": "/maas/api/2.0/nodes/4y3haf/",
        "ip_addresses": ["192.168.100.11"],
        "Owner": "thumper",
        "tag_names": [],
        "FQDN": "furnacelike-brittney.maas",
        "system_id": "4y3haf",
        "Parent": "4y3ha3",
        "interface_set": [
            {
                "resource_uri": "/maas/api/2.0/nodes/4y3haf/interfaces/48/",
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
                "resource_uri": "/maas/api/2.0/nodes/4y3haf/interfaces/49/",
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
	nodesResponse = "[" + nodeResponse + "]"
)
