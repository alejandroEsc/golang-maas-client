// Copyright 2016 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE File for details.

package v2

type Filesystem struct {
	Type       string `json:"Type,omitempty"`
	MountPoint string `json:"mount_point,omitempty"`
	Label      string `json:"Label,omitempty"`
	UUID       string `json:"UUID,omitempty"`
	// no idea what the mount_options are as a value type, so ignoring for now.
}
