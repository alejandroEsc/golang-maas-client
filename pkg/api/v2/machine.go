// Copyright 2016 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE File for details.

package maasapiv2

import (
	"net/http"
	"net/url"

	"encoding/json"

	"github.com/juju/errors"
	"github.com/juju/gomaasapi/pkg/api/client"
	"github.com/juju/gomaasapi/pkg/api/util"
)

// MachineInterface represents a physical MachineInterface.
type Machine struct {
	Controller *Controller `json:"-"`

	ResourceURI string   `json:"resource_uri,omitempty"`
	SystemID    string   `json:"system_id,omitempty"`
	Hostname    string   `json:"Hostname,omitempty"`
	FQDN        string   `json:"FQDN,omitempty"`
	Tags        []string `json:"tag_names,omitempty"`
	// OwnerData returns a copy of the key/value data stored for this
	// object.
	OwnerData       map[string]string `json:"owner_data,omitempty"`
	OperatingSystem string            `json:"osystem,omitempty"`
	DistroSeries    string            `json:"distro_series,omitempty"`
	Architecture    string            `json:"Architecture,omitempty"`
	Memory          int               `json:"Memory,omitempty"`
	CPUCount        int               `json:"cpu_count,omitempty"`
	IPAddresses     []string          `json:"ip_addresses,omitempty"`
	PowerState      string            `json:"power_state,omitempty"`
	// NOTE: consider some form of status struct
	StatusName    string `json:"status_name,omitempty"`
	StatusMessage string `json:"status_message,omitempty"`
	// BootInterface returns the interface that was used to boot the MachineInterface.
	BootInterface *NetworkInterface `json:"boot_interface,omitempty"`
	// InterfaceSet returns all the interfaces for the MachineInterface.
	InterfaceSet []*NetworkInterface `json:"interface_set,omitempty"`
	Zone         *Zone               `json:"Zone,omitempty"`
	// Don't really know the difference between these two lists:

	// PhysicalBlockDevice returns the physical block node for the MachineInterface
	// that matches the ID specified. If there is no match, nil is returned.
	PhysicalBlockDevices []*BlockDevice `json:"physicalblockdevice_set,omitempty"`
	// BlockDevices returns all the physical and virtual block devices on the MachineInterface.
	BlockDevices []*BlockDevice `json:"blockdevice_set,omitempty"`
}

func (m *Machine) updateFrom(other *Machine) {
	m.ResourceURI = other.ResourceURI
	m.SystemID = other.SystemID
	m.Hostname = other.Hostname
	m.FQDN = other.FQDN
	m.OperatingSystem = other.OperatingSystem
	m.DistroSeries = other.DistroSeries
	m.Architecture = other.Architecture
	m.Memory = other.Memory
	m.CPUCount = other.CPUCount
	m.IPAddresses = other.IPAddresses
	m.PowerState = other.PowerState
	m.StatusName = other.StatusName
	m.StatusMessage = other.StatusMessage
	m.Zone = other.Zone
	m.Tags = other.Tags
	m.OwnerData = other.OwnerData
}

// NetworkInterface implements Machine.
func (m *Machine) Interface(id int) *NetworkInterface {
	for _, iface := range m.InterfaceSet {
		if iface.ID == id {
			iface.Controller = m.Controller
			return iface
		}
	}
	return nil
}

// PhysicalBlockDevice implements Machine.
func (m *Machine) PhysicalBlockDevice(id int) *BlockDevice {
	for _, blockDevice := range m.PhysicalBlockDevices {
		if blockDevice.ID == id {
			return blockDevice
		}
	}
	return nil
}

// BlockDevice implements Machine.
func (m *Machine) BlockDevice(id int) *BlockDevice {
	for _, blockDevice := range m.BlockDevices {
		if blockDevice.ID == id {
			return blockDevice
		}
	}
	return nil
}

// Nodes implements Machine.
func (m *Machine) Nodes(args NodesArgs) ([]Node, error) {
	// Perhaps in the future, maas will give us a way to query just for the
	// nodes for a particular Parent.
	nodes, err := m.Controller.Nodes(args)
	if err != nil {
		return nil, errors.Trace(err)
	}

	result := make([]Node, 0)
	for _, n := range nodes {
		if n.Parent == m.SystemID {
			result = append(result, n)
		}
	}
	return result, nil
}

// Deploy implements Machine.
func (m *Machine) Deploy(args DeployMachineArgs) error {
	params := DeploytMachineParams(args)
	result, err := m.Controller.Post(m.ResourceURI, "deploy", params.Values)
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

	machine.Controller = m.Controller

	m.updateFrom(machine)
	return nil
}

// CreateNode implements Machine
func (m *Machine) CreateNode(args CreateMachineNodeArgs) (*Node, error) {
	if err := args.Validate(); err != nil {
		return nil, errors.Trace(err)
	}
	d, err := m.Controller.CreateNode(CreateNodeArgs{
		Hostname:     args.Hostname,
		MACAddresses: []string{args.MACAddress},
		Parent:       m.SystemID,
	})
	if err != nil {
		return nil, err
	}

	defer func(err *error) {
		// If there is an error return, at least try to delete the node we just created.
		if err != nil {
			if innerErr := d.Delete(); innerErr != nil {
				logger.Warningf("could not delete node %q", d.SystemID)
			}
		}
	}(&err)

	// Update the VLAN to use for the interface, if given.
	vlanToUse := args.VLAN
	if vlanToUse == nil && args.Subnet != nil {
		vlanToUse = args.Subnet.VLAN
	}

	// There should be one interface created for each MAC Address, and since we
	// only specified one, there should just be one response.
	interfaces := d.InterfaceSet
	if count := len(interfaces); count != 1 {
		err := errors.Errorf("unexpected interface count for node: %d", count)
		return nil, util.NewUnexpectedError(err)
	}

	if err := m.updateDeviceInterface(interfaces, args.InterfaceName, vlanToUse); err != nil {
		return nil, err
	}

	if args.Subnet == nil {
		return d, nil
	}

	if err := m.linkDeviceInterfaceToSubnet(interfaces, args.Subnet); err != nil {
		return nil, err
	}

	return d, nil
}

func (m *Machine) updateDeviceInterface(interfaces []*NetworkInterface, nameToUse string, vlanToUse *VLAN) error {
	iface := interfaces[0]

	updateArgs := UpdateInterfaceArgs{}
	updateArgs.Name = nameToUse

	if vlanToUse != nil {
		updateArgs.VLAN = vlanToUse
	}

	if err := iface.Update(updateArgs); err != nil {
		return errors.Annotatef(err, "updating node interface %q failed", iface.Name)
	}

	return nil
}

func (m *Machine) linkDeviceInterfaceToSubnet(interfaces []*NetworkInterface, subnetToUse *Subnet) error {
	iface := interfaces[0]

	err := iface.LinkSubnet(LinkSubnetArgs{
		Mode:   LinkModeStatic,
		Subnet: subnetToUse,
	})
	if err != nil {
		return errors.Annotatef(
			err, "linking node interface %q to Subnet %q failed",
			iface.Name, subnetToUse.CIDR)
	}

	return nil
}

// SetOwnerData updates the key/value data stored for this object
// with the Values passed in. Existing keys that aren't specified
// in the map passed in will be left in place; to clear a key set
// its value to "". All Owner data is cleared when the object is
// released.
func (m *Machine) SetOwnerData(ownerData map[string]string) error {
	params := make(url.Values)
	for key, value := range ownerData {
		params.Add(key, value)
	}
	result, err := m.Controller.Post(m.ResourceURI, "set_owner_data", params)
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

func unmarshalMachines(obj []byte) {
}

type MachineInterface interface {
	OwnerDataHolderInterface

	// Nodes returns a list of devices that match the params and have
	// this MachineInterface as the Parent.
	Nodes(NodesArgs) ([]NodeInterface, error)

	// NetworkInterface returns the interface for the MachineInterface that matches the ID
	// specified. If there is no match, nil is returned.
	Interface(id int) *NetworkInterface
	// BlockDevice returns the block node for the MachineInterface that matches the
	// ID specified. If there is no match, nil is returned.
	BlockDevice(id int) BlockDevice

	// Deploy the MachineInterface and install the operating system specified in the args.
	Deploy(DeployMachineArgs) error

	// CreateNode creates a new NodeInterface with this MachineInterface as the Parent.
	// The node will have one interface that is linked to the specified Subnet.
	CreateNode(CreateMachineNodeArgs) (NodeInterface, error)
}

// OwnerDataHolderInterface represents any maas object that can store key/value
// data.
type OwnerDataHolderInterface interface {
	// SetOwnerData updates the key/value data stored for this object
	// with the Values passed in. Existing keys that aren't specified
	// in the map passed in will be left in place; to clear a key set
	// its value to "". All Owner data is cleared when the object is
	// released.
	SetOwnerData(map[string]string) error
}