// Copyright 2016 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE File for details.

package v2

import (
	"fmt"
	"net/http"
	"strings"

	"encoding/json"

	"github.com/alejandroEsc/golang-maas-client/pkg/api/client"
	"github.com/alejandroEsc/golang-maas-client/pkg/api/util"
	"github.com/juju/errors"
)

// Nodes are now known as nodes...? should reconsider this struct.
// Device represents some form of Node in maas.
type Node struct {
	// TODO: add domain
	Controller  *Controller `json:"-"`
	ResourceURI string      `json:"resource_uri,omitempty"`
	SystemID    string      `json:"system_id,omitempty"`
	Hostname    string      `json:"Hostname,omitempty"`
	FQDN        string      `json:"FQDN,omitempty"`
	// Parent returns the SystemID of the Parent. Most often this will be a
	// MachineInterface.
	Parent string `json:"Parent,omitempty"`
	// Owner is the username of the user that created the Node.
	Owner       string   `json:"Owner,omitempty"`
	IPAddresses []string `json:"ip_addresses,omitempty"`
	// InterfaceSet returns all the interfaces for the NodeInterface.
	InterfaceSet []*NetworkInterface `json:"interface_set,omitempty"`
	Zone         *Zone               `json:"Zone,omitempty"`
	Tags         []string            `json:"tag_names,omitempty"`
	Type         string              `json:"node_type_name,omitempty"`
}

// CreateInterfaceArgs is an argument struct for passing parameters to
// the MachineInterface.CreateInterface method.
type CreateInterfaceArgs struct {
	// Name of the interface (required).
	Name string
	// MACAddress is the MAC address of the interface (required).
	MACAddress string
	// VLAN is the untagged VLAN the interface is connected to (required).
	VLAN *VLAN
	// Tags to attach to the interface (optional).
	Tags []string
	// MTU - Maximum transmission unit. (optional)
	MTU int
	// AcceptRA - Accept router advertisements. (IPv6 only)
	AcceptRA bool
	// Autoconf - Perform stateless autoconfiguration. (IPv6 only)
	Autoconf bool
}

// Validate checks the required fields are set for the arg structure.
func (a *CreateInterfaceArgs) Validate() error {
	if a.Name == "" {
		return errors.NotValidf("missing Name")
	}
	if a.MACAddress == "" {
		return errors.NotValidf("missing MACAddress")
	}
	if a.VLAN == nil {
		return errors.NotValidf("missing VLAN")
	}
	return nil
}

// interfacesURI used to add interfaces for this Node. The operations
// are on the nodes endpoint, not devices.
func (d *Node) interfacesURI() string {
	return strings.Replace(d.ResourceURI, "devices", "nodes", 1) + "interfaces/"
}

// CreateInterface implements NodeInterface.
func (d *Node) CreateInterface(args CreateInterfaceArgs) (*NetworkInterface, error) {
	if err := args.Validate(); err != nil {
		return nil, err
	}
	params := util.NewURLParams()
	params.Values.Add("Name", args.Name)
	params.Values.Add("mac_address", args.MACAddress)
	params.Values.Add("VLAN", fmt.Sprint(args.VLAN.ID))
	params.MaybeAdd("Tags", strings.Join(args.Tags, ","))
	params.MaybeAddInt("MTU", args.MTU)
	params.MaybeAddBool("accept_ra", args.AcceptRA)
	params.MaybeAddBool("autoconf", args.Autoconf)

	uri := d.interfacesURI()
	result, err := d.Controller.Post(uri, "create_physical", params.Values)
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
	iface.Controller = d.Controller
	d.InterfaceSet = append(d.InterfaceSet, &iface)
	return &iface, nil
}

// Delete implements NodeInterface.
func (d *Node) Delete() error {
	err := d.Controller.Delete(d.ResourceURI)
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
	return nil
}

type NodeInterface interface {
	// CreateInterface will create a physical interface for this MachineInterface.
	CreateInterface(CreateInterfaceArgs) (*NetworkInterface, error)
	// Delete will remove this NodeInterface.
	Delete() error
}
