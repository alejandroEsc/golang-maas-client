package v2

import (
	"net/http"
	"testing"

	"github.com/alejandroEsc/golang-maas-client/pkg/api/client"
	"github.com/alejandroEsc/golang-maas-client/pkg/api/util"
	"github.com/stretchr/testify/assert"
)

func TestMachineDeploy(t *testing.T) {
	server, machine, controller := getMachineControllerAndServer(t)
	defer server.Close()
	response := util.UpdateJSONMap(t, machineResponse, map[string]interface{}{
		"status_name":    "Deploying",
		"status_message": "for testing",
	})
	server.AddPostResponse(machine.ResourceURI+"?op=deploy", http.StatusOK, response)

	err := controller.Deploy(machine, DeployMachineArgs{
		UserData:     "userdata",
		DistroSeries: "trusty",
		Kernel:       "kernel",
		Comment:      "a comment",
	})
	assert.Nil(t, err)
	assert.Equal(t, machine.StatusName, "Deploying")
	assert.Equal(t, machine.StatusMessage, "for testing")

	request := server.LastRequest()
	// There should be one entry in the form Values for each of the args.
	form := request.PostForm
	assert.Len(t, form, 4)
	assert.Equal(t, form.Get("user_data"), "userdata")
	assert.Equal(t, form.Get("distro_series"), "trusty")
	assert.Equal(t, form.Get("hwe_kernel"), "kernel")
	assert.Equal(t, form.Get("comment"), "a comment")
}

func TestMachineDeployNotFound(t *testing.T) {
	server, machine, controller := getMachineControllerAndServer(t)
	defer server.Close()
	server.AddPostResponse(machine.ResourceURI+"?op=deploy", http.StatusNotFound, "can't find MachineInterface")
	err = controller.Deploy(machine, DeployMachineArgs{})
	assert.True(t, util.IsBadRequestError(err))
	assert.Equal(t, err.Error(), "can't find MachineInterface")
}

func TestMachineDeployConflict(t *testing.T) {
	server, machine, controller := getMachineControllerAndServer(t)
	defer server.Close()
	server.AddPostResponse(machine.ResourceURI+"?op=deploy", http.StatusConflict, "MachineInterface not allocated")
	err = controller.Deploy(machine, DeployMachineArgs{})
	assert.True(t, util.IsBadRequestError(err))
	assert.Equal(t, err.Error(), "MachineInterface not allocated")
}

func TestMachineDeployForbidden(t *testing.T) {
	server, machine, controller := getMachineControllerAndServer(t)
	defer server.Close()
	server.AddPostResponse(machine.ResourceURI+"?op=deploy", http.StatusForbidden, "MachineInterface not yours")
	err = controller.Deploy(machine, DeployMachineArgs{})
	assert.True(t, util.IsPermissionError(err))
	assert.Equal(t, err.Error(), "MachineInterface not yours")
}

func TestMachineDeployServiceUnavailable(t *testing.T) {
	server, machine, controller := getMachineControllerAndServer(t)
	defer server.Close()
	server.AddPostResponse(machine.ResourceURI+"?op=deploy", http.StatusServiceUnavailable, "no ip addresses available")
	err = controller.Deploy(machine, DeployMachineArgs{})
	assert.True(t, util.IsCannotCompleteError(err))
	assert.Equal(t, err.Error(), "no ip addresses available")
}

func TestMachineDeployMachineUnknown(t *testing.T) {
	server, machine, controller := getMachineControllerAndServer(t)
	defer server.Close()
	server.AddPostResponse(machine.ResourceURI+"?op=deploy", http.StatusMethodNotAllowed, "wat?")
	err = controller.Deploy(machine, DeployMachineArgs{})
	assert.True(t, util.IsUnexpectedError(err))
	assert.Equal(t, err.Error(), "unexpected: ServerError: 405 Method Not Allowed (wat?)")
}

func TestMachineSetOwnerData(t *testing.T) {
	server, machine, controller := getMachineControllerAndServer(t)
	server.AddPostResponse(machine.ResourceURI+"?op=set_owner_data", 200, machineWithOwnerData(`{"returned": "data"}`))
	defer server.Close()
	err := controller.SetOwnerData(machine, map[string]string{
		"draco": "malfoy",
		"empty": "", // Check that empty strings get passed along.
	})
	assert.Nil(t, err)
	assert.EqualValues(t, machine.OwnerData, map[string]string{"returned": "data"})
	form := server.LastRequest().PostForm
	// Looking at the map directly so we can tell the difference
	// between no value and an explicit empty string.
	assert.EqualValues(t, form["draco"], []string{"malfoy"})
	assert.EqualValues(t, form["empty"], []string{""})
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
