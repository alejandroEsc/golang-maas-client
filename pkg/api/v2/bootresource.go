// Copyright 2016 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE File for details.

package v2

import (
	"strings"

	"github.com/juju/utils/set"
)

type BootResource struct {
	ResourceURI  string `json:"resource_uri,omitempty"`
	ID           int    `json:"ID,omitempty"`
	Name         string `json:"Name,omitempty"`
	Type         string `json:"type,omitempty"`
	Architecture string `json:"Architecture,omitempty"`
	SubArches    string `json:"subarches,omitempty"`
	KernelFlavor string `json:"kflavor,omitempty"`
}

// SubArchitectures implements BootResource.
func (b *BootResource) SubArchitectures() set.Strings {
	return set.NewStrings(strings.Split(b.SubArches, ",")...)
}
