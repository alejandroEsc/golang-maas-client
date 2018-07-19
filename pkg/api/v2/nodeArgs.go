package v2

import (
	"fmt"
	"strings"

	"github.com/alejandroEsc/golang-maas-client/pkg/api/util"
)

// CreateNodeArgs is a argument struct for passing information into CreateNode.
type CreateNodeArgs struct {
	Hostname     string
	MACAddresses []string
	Domain       string
	Parent       string
}

// NodesArgs is a argument struct for selecting Nodes.
// Only devices that match the specified criteria are returned.
type NodesArgs struct {
	Hostname     []string
	MACAddresses []string
	SystemIDs    []string
	Domain       string
	Zone         string
	AgentName    string
}

func NodesParams(args NodesArgs) *util.URLParams {
	params := util.NewURLParams()
	params.MaybeAddMany("Hostname", args.Hostname)
	params.MaybeAddMany("mac_address", args.MACAddresses)
	params.MaybeAddMany("id", args.SystemIDs)
	params.MaybeAdd("domain", args.Domain)
	params.MaybeAdd("Zone", args.Zone)
	params.MaybeAdd("agent_name", args.AgentName)
	return params
}

func CreateNodesParams(args CreateNodeArgs) *util.URLParams {
	params := util.NewURLParams()
	params.MaybeAdd("Hostname", args.Hostname)
	params.MaybeAdd("domain", args.Domain)
	params.MaybeAddMany("mac_addresses", args.MACAddresses)
	params.MaybeAdd("Parent", args.Parent)
	return params
}

// CreateInterfaceArgs is an argument struct for passing parameters to
// the Machine.CreateInterface method.
type CreateNodeNetworkInterfaceArgs struct {
	// Name of the interface (required).
	Name string
	// MACAddress is the MAC address of the interface (required).
	MACAddress string
	// VLAN is the untagged VLAN the interface is connected to (required).
	VLAN VLAN
	// Tags to attach to the interface (optional).
	Tags []string
	// MTU - Maximum transmission unit. (optional)
	MTU int
	// AcceptRA - Accept router advertisements. (IPv6 only)
	AcceptRA bool
	// Autoconf - Perform stateless autoconfiguration. (IPv6 only)
	Autoconf bool
}

func CreateNodeNetworkInterfaceParams(args CreateNodeNetworkInterfaceArgs) *util.URLParams {
	params := util.NewURLParams()
	params.Values.Add("name", args.Name)
	params.Values.Add("mac_address", args.MACAddress)
	params.Values.Add("vlan", fmt.Sprint(args.VLAN.ID))
	params.MaybeAdd("tags", strings.Join(args.Tags, ","))
	params.MaybeAddInt("mtu", args.MTU)
	params.MaybeAddBool("accept_ra", args.AcceptRA)
	params.MaybeAddBool("autoconf", args.Autoconf)
	return params
}
