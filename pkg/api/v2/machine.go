// Copyright 2016 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE File for details.

package v2

// MachineInterface represents a physical MachineInterface.
type Machine struct {
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
	Kernel       string         `json:"hwe_kernel,omitempty"`
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
