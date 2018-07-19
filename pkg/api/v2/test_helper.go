package v2

import (
	"net/http"
	"testing"

	"github.com/alejandroEsc/golang-maas-client/pkg/api/client"
	"github.com/stretchr/testify/assert"
)

// createTestServerController creates a ControllerInterface backed on to a test server
// that has sufficient knowledge of versions and users to be able to create a
// valid ControllerInterface.
func createTestServerController(t *testing.T) (*client.SimpleTestServer, *Controller) {
	server := client.NewSimpleServer()
	server.AddGetResponse("/api/2.0/users/?op=whoami", http.StatusOK, `"captain awesome"`)
	server.AddGetResponse("/api/2.0/version/", http.StatusOK, versionResponse)
	server.Start()

	controller, err := NewController(ControllerArgs{
		BaseURL: server.URL,
		APIKey:  "fake:as:key",
	})
	assert.Nil(t, err)
	return server, controller
}

func getServeNodeAndController(t *testing.T) (*client.SimpleTestServer, *Node, *Controller) {
	server, controller := createTestServerController(t)
	server.AddGetResponse("/api/2.0/nodes/", http.StatusOK, nodesResponse)

	devices, err := controller.Nodes(NodesArgs{})
	assert.Nil(t, err)
	assert.Len(t, devices, 1)
	return server, &devices[0], controller
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

func getMachineControllerAndServer(t *testing.T) (*client.SimpleTestServer, *Machine, *Controller) {
	server, controller := createTestServerController(t)
	server.AddGetResponse("/api/2.0/machines/", http.StatusOK, "["+machineResponse+"]")

	machines, err := controller.Machines(MachinesArgs{})
	assert.Nil(t, err)
	assert.Len(t, machines, 1)
	machine := machines[0]
	return server, &machine, controller
}

func getController(t *testing.T, server *client.SimpleTestServer) *Controller {
	controller, err := NewController(ControllerArgs{
		BaseURL: server.URL,
		APIKey:  "fake:as:key",
	})
	assert.Nil(t, err)
	return controller
}

func minimalCreateInterfaceArgs() CreateNodeNetworkInterfaceArgs {
	return CreateNodeNetworkInterfaceArgs{
		Name:       "eth43",
		MACAddress: "some-mac-address",
		VLAN:       VLAN{ID: 33},
	}
}
