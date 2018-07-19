package v2

// UpdateInterfaceArgs is an argument struct for calling NetworkInterface.Update.
type UpdateInterfaceArgs struct {
	Name       string
	MACAddress string
	VLAN       *VLAN
}

func (a *UpdateInterfaceArgs) vlanID() int {
	if a.VLAN == nil {
		return 0
	}
	return a.VLAN.ID
}
