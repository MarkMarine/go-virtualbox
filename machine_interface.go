package virtualbox

type Machine interface {
	Refresh() error
	Start() error
	Save() error
	Pause() error
	Stop() error
	Poweroff() error
	Restart() error
	Reset() error
	Delete() error
	Modify() error
	AddNATPF(n int, name string, rule PFRule) error
	DelNATPF(n int, name string) error
	SetNIC(n int, nic NIC) error
	AddStorageCtl(name string, ctl StorageController) error
	DelStorageCtl(name string) error
	AttachStorage(ctlName string, medium StorageMedium) error

	// Getters and Setters
	Name() string
	UUID() string
	State() MachineState
	CPUs() uint
	Memory() uint
	VRAM() uint
	CfgFile() string
	BaseFolder() string
	OSType() string
	Flag() Flag
	BootOrder() []string

	SetName(string)
	SetUUID(string)
	SetState(MachineState)
	SetCPUs(uint)
	SetMemory(uint)
	SetVRAM(uint)
	SetCfgFile(string)
	SetBaseFolder(string)
	SetOSType(string)
	SetFlag(Flag)
	SetBootOrder([]string)
}

// ListMachines lists all registered machines.
func ListMachines() ([]*Machine, error) {
	var mi []*Machine
	machines, err := listMachines()
	if err != nil {
		return mi, err
	}
	for _, machine := range machines {
		i := Machine(machine)
		mi = append(mi, &i)
	}
	return mi, nil
}

// GetMachine finds a machine by its name or UUID.
func GetMachines(id string) (*Machine, error) {
	m, err := getMachine(id)
	if err != nil {
		var x Machine
		return &x, err
	}
	mi := Machine(m)
	return &mi, nil
}

// CreateMachine creates a new machine. If basefolder is empty, use default.
func CreateMachine(name, basefolder string) (*Machine, error) {
	m, err := createMachine(name, basefolder)
	if err != nil {
		var x Machine
		return &x, err
	}
	mi := Machine(m)
	return &mi, nil
}
