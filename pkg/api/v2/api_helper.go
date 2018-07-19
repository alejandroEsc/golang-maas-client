package v2

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"fmt"

	"github.com/alejandroEsc/golang-maas-client/pkg/api/client"
	"github.com/alejandroEsc/golang-maas-client/pkg/api/util"
	"github.com/juju/errors"
)

// GetFile returns a single File by its Filename.
func (c *Controller) GetFile(filename string) (*File, error) {
	if filename == "" {
		return nil, errors.NotValidf("missing Filename")
	}
	source, err := c.Get("files/"+filename, "", nil)
	if err != nil {
		if svrErr, ok := errors.Cause(err).(client.ServerError); ok {
			if svrErr.StatusCode == http.StatusNotFound {
				return nil, errors.Wrap(err, util.NewNoMatchError(svrErr.BodyMessage))
			}
		}
		return nil, util.NewUnexpectedError(err)
	}
	var file File
	err = json.Unmarshal(source, &file)
	if err != nil {
		return nil, err
	}
	return &file, nil
}

func (c *Controller) ReadFileContent(f *File) ([]byte, error) {
	// If the Content is available, it is base64 encoded, so
	args := make(url.Values)
	args.Add("Filename", f.Filename)
	bytes, err := c.Get("files", "Get", args)
	if err != nil {
		if svrErr, ok := errors.Cause(err).(client.ServerError); ok {
			switch svrErr.StatusCode {
			case http.StatusNotFound:
				return nil, errors.Wrap(err, util.NewNoMatchError(svrErr.BodyMessage))
			case http.StatusForbidden:
				return nil, errors.Wrap(err, util.NewPermissionError(svrErr.BodyMessage))
			}
		}
		return nil, util.NewUnexpectedError(err)
	}
	return bytes, nil
}

// getFiles returns all the files that match the specified prefix.
func (c *Controller) getFiles(prefix string) ([]File, error) {
	params := util.NewURLParams()
	params.MaybeAdd("prefix", prefix)
	source, err := c.Get("files", "", params.Values)
	if err != nil {
		return nil, util.NewUnexpectedError(err)
	}

	var files []File
	err = json.Unmarshal(source, &files)
	if err != nil {
		return nil, err
	}
	return files, nil
}

// Fabrics returns the list of Fabrics defined in the maas ControllerInterface.
func (c *Controller) Fabrics() ([]Fabric, error) {
	source, err := c.Get("fabrics", "", nil)
	if err != nil {
		return nil, util.NewUnexpectedError(err)
	}

	var fabrics []Fabric
	err = json.Unmarshal(source, &fabrics)
	if err != nil {
		return nil, err
	}

	return fabrics, nil
}

// Spaces returns the list of Spaces defined in the maas ControllerInterface.
func (c *Controller) Spaces() ([]Space, error) {
	source, err := c.Get("spaces", "", nil)
	if err != nil {
		return nil, util.NewUnexpectedError(err)
	}

	var spaces []Space
	err = json.Unmarshal(source, &spaces)
	if err != nil {
		return nil, err
	}

	return spaces, nil
}

// StaticRoutes returns the list of StaticRoutes defined in the maas ControllerInterface.
func (c *Controller) StaticRoutes() ([]StaticRoute, error) {
	source, err := c.Get("static-routes", "", nil)
	if err != nil {
		return nil, util.NewUnexpectedError(err)
	}
	var staticRoutes []StaticRoute
	err = json.Unmarshal(source, &staticRoutes)
	if err != nil {
		return nil, err
	}

	return staticRoutes, nil
}

// Zones lists all the zones known to the maas ControllerInterface.
func (c *Controller) Zones() ([]Zone, error) {
	source, err := c.Get("zones", "", nil)
	if err != nil {
		return nil, util.NewUnexpectedError(err)
	}
	var zones []Zone
	err = json.Unmarshal(source, &zones)
	if err != nil {
		return nil, err
	}
	return zones, nil
}

// Nodes returns a list of devices that match the params.
func (c *Controller) Nodes(args NodesArgs) ([]Node, error) {
	params := NodesParams(args)
	source, err := c.Get("nodes", "", params.Values)
	if err != nil {
		return nil, util.NewUnexpectedError(err)
	}

	var nodes []Node
	err = json.Unmarshal(source, &nodes)
	if err != nil {
		return nil, err
	}

	return nodes, nil
}

// CreateNode creates and returns a new NodeInterface.
func (c *Controller) CreateNode(args CreateNodeArgs) (*Node, error) {
	// There must be at least one mac address.
	if len(args.MACAddresses) == 0 {
		return nil, util.NewBadRequestError("at least one MAC address must be specified")
	}
	params := CreateNodesParams(args)
	source, err := c.Post("nodes", "", params.Values)
	if err != nil {
		if svrErr, ok := errors.Cause(err).(client.ServerError); ok {
			if svrErr.StatusCode == http.StatusBadRequest {
				return nil, errors.Wrap(err, util.NewBadRequestError(svrErr.BodyMessage))
			}
		}
		// Translate http errors.
		return nil, util.NewUnexpectedError(err)
	}

	var d Node
	err = json.Unmarshal(source, &d)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

// Machines returns a list of machines that match the params.
func (c *Controller) Machines(args MachinesArgs) ([]Machine, error) {
	params := MachinesParams(args)
	source, err := c.Get("machines", "", params.Values)
	if err != nil {
		return nil, util.NewUnexpectedError(err)
	}

	var machines []Machine
	err = json.Unmarshal(source, &machines)
	if err != nil {
		return nil, err
	}

	return machines, nil
}

// AddFile adds or replaces the Content of the specified Filename.
// If or when the maas api is able to return metadata about a single
// File without sending the Content of the File, we can return a FileInterface
// instance here too.
func (c *Controller) AddFile(args AddFileArgs) error {
	if err := args.Validate(); err != nil {
		return err
	}
	fileContent := args.Content
	if fileContent == nil {
		content, err := ioutil.ReadAll(io.LimitReader(args.Reader, args.Length))
		if err != nil {
			return errors.Annotatef(err, "cannot read File Content")
		}
		fileContent = content
	}
	params := url.Values{"Filename": {args.Filename}}
	_, err := c.PostFile("files", "", params, fileContent)
	if err != nil {
		if svrErr, ok := errors.Cause(err).(client.ServerError); ok {
			if svrErr.StatusCode == http.StatusBadRequest {
				return errors.Wrap(err, util.NewBadRequestError(svrErr.BodyMessage))
			}
		}
		return util.NewUnexpectedError(err)
	}
	return nil
}

func ownerDataMatches(ownerData, filter map[string]string) bool {
	for key, value := range filter {
		if ownerData[key] != value {
			return false
		}
	}
	return true
}

// BootResources implements ControllerInterface.
func (c *Controller) BootResources() ([]*BootResource, error) {
	source, err := c.Get("boot-resources", "", nil)
	if err != nil {
		return nil, util.NewUnexpectedError(err)
	}

	var resources []*BootResource
	err = json.Unmarshal(source, &resources)
	if err != nil {
		return nil, err
	}

	return resources, nil
}

// AllocateMachine will attempt to allocate a MachineInterface to the user.
// If successful, the allocated MachineInterface is returned.
// Returns an error that satisfies IsNoMatchError if the requested
// constraints cannot be met.
func (c *Controller) AllocateMachine(args AllocateMachineArgs) (*Machine, ConstraintMatches, error) {
	var matches ConstraintMatches
	params := AllocateMachinesParams(args)
	result, err := c.Post("machines", "allocate", params.Values)
	if err != nil {
		// A 409 Status code is "No Matching Machines"
		if svrErr, ok := errors.Cause(err).(client.ServerError); ok {
			if svrErr.StatusCode == http.StatusConflict {
				return nil, matches, errors.Wrap(err, util.NewNoMatchError(svrErr.BodyMessage))
			}
		}
		// Translate http errors.
		return nil, matches, util.NewUnexpectedError(err)
	}

	var machine *Machine
	var source map[string]interface{}
	err = json.Unmarshal(result, &machine)
	if err != nil {
		return nil, matches, err
	}
	err = json.Unmarshal(result, &source)
	if err != nil {
		return nil, matches, err
	}

	// Parse the constraint matches.
	matches, err = parseAllocateConstraintsResponse(source, machine)
	if err != nil {
		return nil, matches, err
	}

	return machine, matches, nil
}

// ReleaseMachines will stop the specified machines, and release them
// from the user making them available to be allocated again.
// Release multiple machines at once. Returns
//  - BadRequestError if any of the machines cannot be found
//  - PermissionError if the user does not have permission to release any of the machines
//  - CannotCompleteError if any of the machines could not be released due to their current state
func (c *Controller) ReleaseMachines(args ReleaseMachinesArgs) error {
	params := ReleaseMachinesParams(args)
	_, err := c.Post("machines", "release", params.Values)
	if err != nil {
		if svrErr, ok := errors.Cause(err).(client.ServerError); ok {
			switch svrErr.StatusCode {
			case http.StatusBadRequest:
				return errors.Wrap(err, util.NewBadRequestError(svrErr.BodyMessage))
			case http.StatusForbidden:
				return errors.Wrap(err, util.NewPermissionError(svrErr.BodyMessage))
			case http.StatusConflict:
				return errors.Wrap(err, util.NewCannotCompleteError(svrErr.BodyMessage))
			}
		}
		return util.NewUnexpectedError(err)
	}

	return nil
}

// Deploy implements Machine.
func (c *Controller) Deploy(m *Machine, args DeployMachineArgs) error {
	params := DeploytMachineParams(args)
	result, err := c.Post(m.ResourceURI, "deploy", params.Values)
	if err != nil {
		if svrErr, ok := errors.Cause(err).(client.ServerError); ok {
			switch svrErr.StatusCode {
			case http.StatusNotFound, http.StatusConflict:
				return errors.Wrap(err, util.NewBadRequestError(svrErr.BodyMessage))
			case http.StatusForbidden:
				return errors.Wrap(err, util.NewPermissionError(svrErr.BodyMessage))
			case http.StatusServiceUnavailable:
				return errors.Wrap(err, util.NewCannotCompleteError(svrErr.BodyMessage))
			}
		}
		return util.NewUnexpectedError(err)
	}

	var machine *Machine
	err = json.Unmarshal(result, &machine)
	if err != nil {
		return errors.Trace(err)
	}
	m.updateFrom(machine)
	return nil
}

// Devices implements Controller.
func (c *Controller) Devices(args DevicesArgs) ([]Device, error) {
	params := GetDeviceParams(args)
	source, err := c.Get("devices", "", params.Values)
	if err != nil {
		return nil, util.NewUnexpectedError(err)
	}

	devices := make([]Device, 0)
	err = json.Unmarshal(source, &devices)
	if err != nil {
		return nil, err
	}

	return devices, nil
}

// SetOwnerData updates the key/value data stored for this object
// with the Values passed in. Existing keys that aren't specified
// in the map passed in will be left in place; to clear a key set
// its value to "". All Owner data is cleared when the object is
// released.
func (c *Controller) SetOwnerData(m *Machine, ownerData map[string]string) error {
	params := make(url.Values)
	for key, value := range ownerData {
		params.Add(key, value)
	}
	result, err := c.Post(m.ResourceURI, "set_owner_data", params)
	if err != nil {
		return errors.Trace(err)
	}

	var machine *Machine
	err = json.Unmarshal(result, &machine)
	if err != nil {
		return errors.Trace(err)
	}

	m.updateFrom(machine)
	return nil
}

// CreateInterface implements NodeInterface.
func (c *Controller) CreateInterface(d *Node, args CreateNodeNetworkInterfaceArgs) (*NetworkInterface, error) {
	params := CreateNodeNetworkInterfaceParams(args)
	result, err := c.Post(d.ResourceURI+"interfaces/", "create_physical", params.Values)
	if err != nil {
		if svrErr, ok := errors.Cause(err).(client.ServerError); ok {
			switch svrErr.StatusCode {
			case http.StatusNotFound, http.StatusConflict:
				return nil, errors.Wrap(err, util.NewBadRequestError(svrErr.BodyMessage))
			case http.StatusForbidden:
				return nil, errors.Wrap(err, util.NewPermissionError(svrErr.BodyMessage))
			case http.StatusServiceUnavailable:
				return nil, errors.Wrap(err, util.NewCannotCompleteError(svrErr.BodyMessage))
			}
		}
		return nil, util.NewUnexpectedError(err)
	}

	var iface NetworkInterface
	err = json.Unmarshal(result, &iface)
	if err != nil {
		return nil, err
	}
	d.InterfaceSet = append(d.InterfaceSet, &iface)
	return &iface, nil
}

// UnlinkSubnet will remove the Link to the Subnet, and release the IP
// address associated if there is one.
func (c *Controller) UnlinkSubnet(i *NetworkInterface, s *Subnet) error {
	if s == nil {
		return errors.NotValidf("missing Subnet")
	}
	link := i.linkForSubnet(s)
	if link == nil {
		return errors.NotValidf("unlinked Subnet")
	}
	params := util.NewURLParams()
	params.Values.Add("ID", fmt.Sprint(link.ID))
	source, err := c.Post(i.ResourceURI, "unlink_subnet", params.Values)
	if err != nil {
		if svrErr, ok := errors.Cause(err).(client.ServerError); ok {
			switch svrErr.StatusCode {
			case http.StatusNotFound, http.StatusBadRequest:
				return errors.Wrap(err, util.NewBadRequestError(svrErr.BodyMessage))
			case http.StatusForbidden:
				return errors.Wrap(err, util.NewPermissionError(svrErr.BodyMessage))
			}
		}
		return util.NewUnexpectedError(err)
	}

	var response NetworkInterface
	err = json.Unmarshal(source, &response)
	if err != nil {
		return errors.Trace(err)
	}

	i.updateFrom(&response)

	return nil
}

// LinkSubnet will attempt to make this interface available on the specified
// Subnet.
func (c *Controller) LinkSubnet(i *NetworkInterface, args LinkSubnetArgs) error {
	if err := args.Validate(); err != nil {
		return errors.Trace(err)
	}
	params := util.NewURLParams()
	params.Values.Add("Mode", string(args.Mode))
	params.Values.Add("Subnet", fmt.Sprint(args.Subnet.ID))
	params.MaybeAdd("ip_address", args.IPAddress)
	params.MaybeAddBool("default_gateway", args.DefaultGateway)
	source, err := c.Post(i.ResourceURI, "link_subnet", params.Values)
	if err != nil {
		if svrErr, ok := errors.Cause(err).(client.ServerError); ok {
			switch svrErr.StatusCode {
			case http.StatusNotFound, http.StatusBadRequest:
				return errors.Wrap(err, util.NewBadRequestError(svrErr.BodyMessage))
			case http.StatusForbidden:
				return errors.Wrap(err, util.NewPermissionError(svrErr.BodyMessage))
			case http.StatusServiceUnavailable:
				return errors.Wrap(err, util.NewCannotCompleteError(svrErr.BodyMessage))
			}
		}
		return util.NewUnexpectedError(err)
	}

	var response NetworkInterface
	err = json.Unmarshal(source, &response)
	if err != nil {
		return errors.Trace(err)
	}

	i.updateFrom(&response)
	return nil
}

// Update the Name, mac address or VLAN.
func (c *Controller) UpdateNetworkInterface(i *NetworkInterface, args UpdateInterfaceArgs) error {
	var empty UpdateInterfaceArgs
	if args == empty {
		return fmt.Errorf("params are empty, and are required.")
	}

	params := UpdateInterfaceParams(args)

	source, err := c.Put(i.ResourceURI, params.Values)
	if err != nil {
		if svrErr, ok := errors.Cause(err).(client.ServerError); ok {
			switch svrErr.StatusCode {
			case http.StatusNotFound:
				return errors.Wrap(err, util.NewNoMatchError(svrErr.BodyMessage))
			case http.StatusForbidden:
				return errors.Wrap(err, util.NewPermissionError(svrErr.BodyMessage))
			}
		}
		return util.NewUnexpectedError(err)
	}

	var response NetworkInterface
	err = json.Unmarshal(source, &response)
	if err != nil {
		return errors.Trace(err)
	}
	i.updateFrom(&response)
	return nil
}

func (c *Controller) linkDeviceInterfaceToSubnet(m *Machine, interfaces []*NetworkInterface, subnetToUse *Subnet) error {
	iface := interfaces[0]
	args := LinkSubnetArgs{
		Mode:   LinkModeStatic,
		Subnet: subnetToUse,
	}

	err := c.LinkSubnet(iface, args)
	if err != nil {
		return errors.Annotatef(
			err, "linking node interface %q to Subnet %q failed",
			iface.Name, subnetToUse.CIDR)
	}

	return nil
}

func (c *Controller) updateDeviceInterface(interfaces []*NetworkInterface, nameToUse string, vlanToUse *VLAN) error {
	iface := interfaces[0]

	args := UpdateInterfaceArgs{Name: nameToUse}
	if vlanToUse != nil {
		args.VLAN = *vlanToUse
	}

	if err := c.UpdateNetworkInterface(iface, args); err != nil {
		return errors.Annotatef(err, "updating node interface %q failed", iface.Name)
	}

	return nil
}
