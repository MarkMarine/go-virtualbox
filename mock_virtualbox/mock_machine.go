package mock_virtualbox

import (
	"github.com/markmarine/go-virtualbox"
	"github.com/satori/go.uuid"
)

// Machine information.
type MockMachine struct {
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

func (m *MockMachine) Refresh() error {
	return nil
}

func (m *MockMachine) Start() error {
	m.state = virtualbox.Running
	return nil
}

func (m *MockMachine) Save() error {
	m.state = virtualbox.Saved
	return nil
}

func (m *MockMachine) Pause() error {
	m.state = virtualbox.Paused
	return nil
}

// Stop gracefully stops the machine.
func (m *MockMachine) Stop() error {
	m.state = virtualbox.Poweroff
	return nil
}

// Poweroff forcefully stops the machine. State is lost and might corrupt the disk image.
func (m *MockMachine) Poweroff() error {
	m.state = virtualbox.Poweroff
	return nil
}

// Restart gracefully restarts the machine.
func (m *MockMachine) Restart() error {
	return m.Start()
}

// Reset forcefully restarts the machine. State is lost and might corrupt the disk image.
func (m *MockMachine) Reset() error {
	return m.Start()
}

// Delete deletes the machine and associated disk images.
func (m *MockMachine) Delete() error {
	x := MockMachine{}
	m = &x
	return nil
}

// GetMachine finds a machine by its name or UUID.
func GetMachine(id string) (*MockMachine, error) {
	// TODO find a better way to do this
	m := MockMachine{name: id}
	return &m, nil
}

// ListMachines lists all registered machines.
func ListMachines() ([]*MockMachine, error) {
	return ms, nil
}

// CreateMachine creates a new machine. If basefolder is empty, use default.
func CreateMachine(name, basefolder string) (*MockMachine, error) {
	var m MockMachine
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
func (m *MockMachine) Modify() error {
	return m.Refresh()
}

// AddNATPF adds a NAT port forarding rule to the n-th NIC with the given name.
func (m *MockMachine) AddNATPF(n int, name string, rule virtualbox.PFRule) error {
	return nil
}

// DelNATPF deletes the NAT port forwarding rule with the given name from the n-th NIC.
func (m *MockMachine) DelNATPF(n int, name string) error {
	return nil
}

// SetNIC set the n-th NIC.
func (m *MockMachine) SetNIC(n int, nic virtualbox.NIC) error {
	return nil
}

// AddStorageCtl adds a storage controller with the given name.
func (m *MockMachine) AddStorageCtl(name string, ctl virtualbox.StorageController) error {
	return nil
}

// DelStorageCtl deletes the storage controller with the given name.
func (m *MockMachine) DelStorageCtl(name string) error {
	return nil
}

// AttachStorage attaches a storage medium to the named storage controller.
func (m *MockMachine) AttachStorage(ctlName string, medium virtualbox.StorageMedium) error {
	return nil
}

func (m *MockMachine) Name() string {
	return m.name
}

func (m *MockMachine) UUID() string {
	return m.uUID
}

func (m *MockMachine) State() virtualbox.MachineState {
	return m.state
}

func (m *MockMachine) CPUs() uint {
	return m.cPUs
}

func (m *MockMachine) Memory() uint {
	return m.memory
}

func (m *MockMachine) VRAM() uint {
	return m.vRAM
}

func (m *MockMachine) CfgFile() string {
	return m.cfgFile
}

func (m *MockMachine) BaseFolder() string {
	return m.baseFolder
}

func (m *MockMachine) OSType() string {
	return m.oSType
}

func (m *MockMachine) Flag() virtualbox.Flag {
	return m.flag
}

func (m *MockMachine) BootOrder() []string {
	return m.bootOrder
}

func (m *MockMachine) SetName(name string) {
	m.name = name
}

func (m *MockMachine) SetUUID(uuid string) {
	// TODO add validation
	m.uUID = uuid
}

func (m *MockMachine) SetState(state virtualbox.MachineState) {
	m.state = state
}

func (m *MockMachine) SetCPUs(cpus uint) {
	m.cPUs = cpus
}

func (m *MockMachine) SetMemory(memory uint) {
	m.memory = memory
}

func (m *MockMachine) SetVRAM(vram uint) {
	m.vRAM = vram
}

func (m *MockMachine) SetCfgFile(cfgFile string) {
	m.cfgFile = cfgFile
}

func (m *MockMachine) SetBaseFolder(baseFolder string) {
	m.baseFolder = baseFolder
}

func (m *MockMachine) SetOSType(osType string) {
	m.oSType = osType
}

func (m *MockMachine) SetFlag(flag virtualbox.Flag) {
	m.flag = flag
}

func (m *MockMachine) SetBootOrder(bootOrder []string) {
	m.bootOrder = bootOrder
}
