package mock_virtualbox

import (
	"github.com/markmarine/go-virtualbox"
	"github.com/satori/go.uuid"
)

// Machine information.
type testMockMachine struct {
	name       string
	uUID       string
	state      virtualbox.MachineState
	cPUs       uint
	memory     uint // main memory (in MB)
	vRAM       uint // video memory (in MB)
	cfgFile    string
	baseFolder string
	oSType     string
	flag       virtualbox.Flag
	bootOrder  []string // max 4 slots, each in {none|floppy|dvd|disk|net}
}

func (m *testMockMachine) Refresh() error {
	return nil
}

func (m *testMockMachine) Start() error {
	m.state = virtualbox.Running
	return nil
}

func (m *testMockMachine) Save() error {
	m.state = virtualbox.Saved
	return nil
}

func (m *testMockMachine) Pause() error {
	m.state = virtualbox.Paused
	return nil
}

// Stop gracefully stops the machine.
func (m *testMockMachine) Stop() error {
	m.state = virtualbox.Poweroff
	return nil
}

// Poweroff forcefully stops the machine. State is lost and might corrupt the disk image.
func (m *testMockMachine) Poweroff() error {
	m.state = virtualbox.Poweroff
	return nil
}

// Restart gracefully restarts the machine.
func (m *testMockMachine) Restart() error {
	return m.Start()
}

// Reset forcefully restarts the machine. State is lost and might corrupt the disk image.
func (m *testMockMachine) Reset() error {
	return m.Start()
}

// Delete deletes the machine and associated disk images.
func (m *testMockMachine) Delete() error {
	x := testMockMachine{}
	m = &x
	return nil
}

// GetMachine finds a machine by its name or UUID.
func GetMachine(id string) (*testMockMachine, error) {
	// TODO find a better way to do this
	m := testMockMachine{name: id}
	return &m, nil
}

// ListMachines lists all registered machines.
func ListMachines() ([]*testMockMachine, error) {
	return ms, nil
}

// CreateMachine creates a new machine. If basefolder is empty, use default.
func CreateMachine(name, basefolder string) (*testMockMachine, error) {
	var m testMockMachine
	if name == "" {
		m.name == "default"
	}
	m.baseFolder = basefolder
	m.uUID = uuid.NewV4()
	m.state = virtualbox.Running
	m.memory = 512
	m.vRAM = 12
	m.cfgFile = ""
	m.baseFolder = ""
	m.oSType = "Ubuntu (64-bit)"
	m.flag = 1
	m.bootOrder = []string{"disk"}
	return &m, nil
}

// Modify changes the settings of the machine.
func (m *testMockMachine) Modify() error {
	return m.Refresh()
}

// AddNATPF adds a NAT port forarding rule to the n-th NIC with the given name.
func (m *testMockMachine) AddNATPF(n int, name string, rule virtualbox.PFRule) error {
	return nil
}

// DelNATPF deletes the NAT port forwarding rule with the given name from the n-th NIC.
func (m *testMockMachine) DelNATPF(n int, name string) error {
	return nil
}

// SetNIC set the n-th NIC.
func (m *testMockMachine) SetNIC(n int, nic virtualbox.NIC) error {
	return nil
}

// AddStorageCtl adds a storage controller with the given name.
func (m *testMockMachine) AddStorageCtl(name string, ctl virtualbox.StorageController) error {
	return nil
}

// DelStorageCtl deletes the storage controller with the given name.
func (m *testMockMachine) DelStorageCtl(name string) error {
	return nil
}

// AttachStorage attaches a storage medium to the named storage controller.
func (m *testMockMachine) AttachStorage(ctlName string, medium virtualbox.StorageMedium) error {
	return nil
}

func (m *testMockMachine) Name() string {
	return m.name
}

func (m *testMockMachine) UUID() string {
	return m.uUID
}

func (m *testMockMachine) State() virtualbox.MachineState {
	return m.state
}

func (m *testMockMachine) CPUs() uint {
	return m.cPUs
}

func (m *testMockMachine) Memory() uint {
	return m.memory
}

func (m *testMockMachine) VRAM() uint {
	return m.vRAM
}

func (m *testMockMachine) CfgFile() string {
	return m.cfgFile
}

func (m *testMockMachine) BaseFolder() string {
	return m.baseFolder
}

func (m *testMockMachine) OSType() string {
	return m.oSType
}

func (m *testMockMachine) Flag() virtualbox.Flag {
	return m.flag
}

func (m *testMockMachine) BootOrder() []string {
	return m.bootOrder
}

func (m *testMockMachine) SetName(name string) {
	m.name = name
}

func (m *testMockMachine) SetUUID(uuid string) {
	// TODO add validation
	m.uUID = uuid
}

func (m *testMockMachine) SetState(state virtualbox.MachineState) {
	m.state = state
}

func (m *testMockMachine) SetCPUs(cpus uint) {
	m.cPUs = cpus
}

func (m *testMockMachine) SetMemory(memory uint) {
	m.memory = memory
}

func (m *testMockMachine) SetVRAM(vram uint) {
	m.vRAM = vram
}

func (m *testMockMachine) SetCfgFile(cfgFile string) {
	m.cfgFile = cfgFile
}

func (m *testMockMachine) SetBaseFolder(baseFolder string) {
	m.baseFolder = baseFolder
}

func (m *testMockMachine) SetOSType(osType string) {
	m.oSType = osType
}

func (m *testMockMachine) SetFlag(flag virtualbox.Flag) {
	m.flag = flag
}

func (m *testMockMachine) SetBootOrder(bootOrder []string) {
	m.bootOrder = bootOrder
}
