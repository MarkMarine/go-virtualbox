package mock_virtualbox

import (
	"github.com/markmarine/go-virtualbox"
	"github.com/satori/go.uuid"
	"fmt"
	"errors"
)

// Machine information.
type MockMachineErr struct {
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

var mockErr error = errors.New("mock os exit 1")

func (m *MockMachineErr) Refresh() error {
	return mockErr
}

func (m *MockMachineErr) Start() error {
	return mockErr
}

func (m *MockMachineErr) Save() error {
	return mockErr
}

func (m *MockMachineErr) Pause() error {
	return mockErr
}

// Stop gracefully stops the machine.
func (m *MockMachineErr) Stop() error {
	return mockErr
}

// Poweroff forcefully stops the machine. State is lost and might corrupt the disk image.
func (m *MockMachineErr) Poweroff() error {
	return mockErr
}

// Restart gracefully restarts the machine.
func (m *MockMachineErr) Restart() error {
	return m.Start()
}

// Reset forcefully restarts the machine. State is lost and might corrupt the disk image.
func (m *MockMachineErr) Reset() error {
	return m.Start()
}

// Delete deletes the machine and associated disk images.
func (m *MockMachineErr) Delete() error {
	return mockErr
}


// CreateMachine creates a new machine. If basefolder is empty, use default.
func CreateMachineErr(name, basefolder string) (*MockMachineErr, error) {
	var m MockMachineErr
	if name == "" {
		m.name = "default"
	} else {
		m.name = name
	}
	m.baseFolder = basefolder
	m.uUID = fmt.Sprintf("%v", uuid.NewV4())
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
func (m *MockMachineErr) Modify() error {
	return m.Refresh()
}

// AddNATPF adds a NAT port forarding rule to the n-th NIC with the given name.
func (m *MockMachineErr) AddNATPF(n int, name string, rule virtualbox.PFRule) error {
	return mockErr
}

// DelNATPF deletes the NAT port forwarding rule with the given name from the n-th NIC.
func (m *MockMachineErr) DelNATPF(n int, name string) error {
	return mockErr
}

// SetNIC set the n-th NIC.
func (m *MockMachineErr) SetNIC(n int, nic virtualbox.NIC) error {
	return mockErr
}

// AddStorageCtl adds a storage controller with the given name.
func (m *MockMachineErr) AddStorageCtl(name string, ctl virtualbox.StorageController) error {
	return mockErr
}

// DelStorageCtl deletes the storage controller with the given name.
func (m *MockMachineErr) DelStorageCtl(name string) error {
	return mockErr
}

// AttachStorage attaches a storage medium to the named storage controller.
func (m *MockMachineErr) AttachStorage(ctlName string, medium virtualbox.StorageMedium) error {
	return mockErr
}

func (m *MockMachineErr) Name() string {
	return m.name
}

func (m *MockMachineErr) UUID() string {
	return m.uUID
}

func (m *MockMachineErr) State() virtualbox.MachineState {
	return m.state
}

func (m *MockMachineErr) CPUs() uint {
	return m.cPUs
}

func (m *MockMachineErr) Memory() uint {
	return m.memory
}

func (m *MockMachineErr) VRAM() uint {
	return m.vRAM
}

func (m *MockMachineErr) CfgFile() string {
	return m.cfgFile
}

func (m *MockMachineErr) BaseFolder() string {
	return m.baseFolder
}

func (m *MockMachineErr) OSType() string {
	return m.oSType
}

func (m *MockMachineErr) Flag() virtualbox.Flag {
	return m.flag
}

func (m *MockMachineErr) BootOrder() []string {
	return m.bootOrder
}

func (m *MockMachineErr) SetName(name string) {
	m.name = name
}

func (m *MockMachineErr) SetUUID(uuid string) {
	m.uUID = uuid
}

func (m *MockMachineErr) SetState(state virtualbox.MachineState) {
	m.state = state
}

func (m *MockMachineErr) SetCPUs(cpus uint) {
	m.cPUs = cpus
}

func (m *MockMachineErr) SetMemory(memory uint) {
	m.memory = memory
}

func (m *MockMachineErr) SetVRAM(vram uint) {
	m.vRAM = vram
}

func (m *MockMachineErr) SetCfgFile(cfgFile string) {
	m.cfgFile = cfgFile
}

func (m *MockMachineErr) SetBaseFolder(baseFolder string) {
	m.baseFolder = baseFolder
}

func (m *MockMachineErr) SetOSType(osType string) {
	m.oSType = osType
}

func (m *MockMachineErr) SetFlag(flag virtualbox.Flag) {
	m.flag = flag
}

func (m *MockMachineErr) SetBootOrder(bootOrder []string) {
	m.bootOrder = bootOrder
}
