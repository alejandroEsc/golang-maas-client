package v2

import "github.com/alejandroEsc/golang-maas-client/pkg/api/util"

// DevicesArgs is a argument struct for selecting Devices.
// Only devices that match the specified criteria are returned.
type DevicesArgs struct {
	Hostname     []string
	MACAddresses []string
	SystemIDs    []string
	Domain       string
	Zone         string
	AgentName    string
}

func GetDeviceParams(args DevicesArgs) *util.URLParams {
	params := util.NewURLParams()
	params.MaybeAddMany("hostname", args.Hostname)
	params.MaybeAddMany("mac_address", args.MACAddresses)
	params.MaybeAddMany("id", args.SystemIDs)
	params.MaybeAdd("domain", args.Domain)
	params.MaybeAdd("zone", args.Zone)
	params.MaybeAdd("agent_name", args.AgentName)
	return params
}
