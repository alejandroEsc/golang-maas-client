package v2

type MachineOp string

const (

	// Comission begins commissioning process for a machine.
	MachineComission MachineOp = "comission"
	// Deploy an operating system to a machine.
	MachineDeploy MachineOp = "deploy"
	// Details obtains various system details.
	MachineDetails MachineOp = "details"
	// GetCurtinConfig returns the rendered curtin configuration for the machine.
	MachineGetCurtinConfig MachineOp = "get_curtin_config"
	// PowerParams obtain power parameters.
	PowerParams MachineOp = "power_parameters"
	// Abort a machine's current operation
	MachineAbort MachineOp = "abort"
	// clear_default_gateways
	MachineClearDefaultGateways MachineOp = "clear_default_gateways"
	// ExitRescueMode exits rescue mode process for a machine.
	MachineExitRescueMode MachineOp = "exit_rescue_mode"
	// MarkBroken marks a node as 'broken'.
	MachineMarkBroken MachineOp = "mark_broken"
	// MarkFixed mark a broken node as fixed and set its status as 'ready'.
	MachineMarkFixed MachineOp = "mark_fixed"
	// MountSpecial Mount a special-purpose filesystem, like tmpfs.
	MachineMountSpecial MachineOp = "mount_special"
	// PowerOFF to request Power off a node.
	MachinePowerOFF MachineOp = "power_off"
	// PowerON Turn on a node.
	MachinePowerON MachineOp = "power_on"
	// Release  a machine. Opposite of Machines.allocate.
	MachineRelease MachineOp = "release"
	// Begin rescue mode process for a machine.
	MachineRescueMode MachineOp = "rescue_mode"
	// Reset a machine's configuration to its initial state.
	MachineRestoreDefaultConfig MachineOp = "restore_default_configuration"
	// Reset a machine's networking options to its initial state.
	MachineRestoreNetworkConfig MachineOp = "restore_networking_configuration"
	// Reset a machine's storage options to its initial state.
	MachineRestoreStorageConfig MachineOp = "restore_storage_configuration"
	// Set key/value data for the current owner.
	MachineSetOwnerData MachineOp = "set_owner_data"
	// Changes the storage layout on the machine.
	MachineSetStorageLayout MachineOp = "set_storage_layout"
	// Unmount a special-purpose filesystem, like tmpfs.
	MachineUnmountSpecial MachineOp = "unmount_special"
)

type MachinesOp string

const (
	// Allocate an available machine for deployment.
	Allocate MachinesOp = "allocate"
)
