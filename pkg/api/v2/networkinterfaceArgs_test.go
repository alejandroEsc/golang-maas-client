package v2

import (
	"testing"

	"github.com/juju/errors"
	"github.com/stretchr/testify/assert"
)

func TestNetworkInterfaceLinkSubnetArgs(t *testing.T) {
	for _, test := range []struct {
		args    LinkSubnetArgs
		errText string
	}{{
		errText: "missing Mode not valid",
	}, {
		args:    LinkSubnetArgs{Mode: LinkModeDHCP},
		errText: "missing Subnet not valid",
	}, {
		args:    LinkSubnetArgs{Mode: InterfaceLinkMode("foo")},
		errText: `unknown Mode value ("foo") not valid`,
	}, {
		args: LinkSubnetArgs{Mode: LinkModeDHCP, Subnet: &Subnet{}},
	}, {
		args: LinkSubnetArgs{Mode: LinkModeStatic, Subnet: &Subnet{}},
	}, {
		args: LinkSubnetArgs{Mode: LinkModeLinkUp, Subnet: &Subnet{}},
	}, {
		args:    LinkSubnetArgs{Mode: LinkModeDHCP, Subnet: &Subnet{}, IPAddress: "10.10.10.10"},
		errText: `setting IP Address when Mode is not LinkModeStatic not valid`,
	}, {
		args: LinkSubnetArgs{Mode: LinkModeStatic, Subnet: &Subnet{}, IPAddress: "10.10.10.10"},
	}, {
		args:    LinkSubnetArgs{Mode: LinkModeLinkUp, Subnet: &Subnet{}, IPAddress: "10.10.10.10"},
		errText: `setting IP Address when Mode is not LinkModeStatic not valid`,
	}, {
		args:    LinkSubnetArgs{Mode: LinkModeDHCP, Subnet: &Subnet{}, DefaultGateway: true},
		errText: `specifying DefaultGateway for Mode "DHCP" not valid`,
	}, {
		args: LinkSubnetArgs{Mode: LinkModeStatic, Subnet: &Subnet{}, DefaultGateway: true},
	}, {
		args:    LinkSubnetArgs{Mode: LinkModeLinkUp, Subnet: &Subnet{}, DefaultGateway: true},
		errText: `specifying DefaultGateway for Mode "LINK_UP" not valid`,
	}} {
		err := test.args.Validate()
		if test.errText == "" {
			assert.Nil(t, err)
		} else {
			assert.True(t, errors.IsNotValid(err))
			assert.Equal(t, err.Error(), test.errText)
		}
	}
}
