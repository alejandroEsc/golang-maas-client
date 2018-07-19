package v2

import "github.com/alejandroEsc/golang-maas-client/pkg/api/util"

// UpdateInterfaceArgs is an argument struct for calling NetworkInterface.Update.
type UpdateInterfaceArgs struct {
	Name       string
	MACAddress string
	VLAN       VLAN
}

func UpdateInterfaceParams(args UpdateInterfaceArgs) *util.URLParams {
	params := util.NewURLParams()
	params.MaybeAdd("Name", args.Name)
	params.MaybeAdd("mac_address", args.MACAddress)
	params.MaybeAddInt("VLAN", args.VLAN.ID)
	return params
}
