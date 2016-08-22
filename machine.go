package virtualbox

import (
	"bufio"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type MachineState string

const (
	Poweroff = MachineState("poweroff")
	Running  = MachineState("running")
	Paused   = MachineState("paused")
	Saved    = MachineState("saved")
	Aborted  = MachineState("aborted")
)

type Flag int

// Flag names in lowercases to be consistent with VBoxManage options.
const (
	F_acpi Flag = 1 << iota
	F_ioapic
	F_rtcuseutc
	F_cpuhotplug
	F_pae
	F_longmode
	F_synthcpu
	F_hpet
	F_hwvirtex
	F_triplefaultreset
	F_nestedpaging
	F_largepages
	F_vtxvpid
	F_vtxux
	F_accelerate3d
)

// Convert bool to "on"/"off"
func bool2string(b bool) string {
	if b {
		return "on"
	}
	return "off"
}

// Test if flag is set. Return "on" or "off".
func (f Flag) Get(o Flag) string {
	return bool2string(f&o == o)
}

// Machine information.
type machine struct {
	name       string
	uUID       string
	state      MachineState
	cPUs       uint
	memory     uint // main memory (in MB)
	vRAM       uint // video memory (in MB)
	cfgFile    string
	baseFolder string
	oSType     string
	flag       Flag
	bootOrder  []string // max 4 slots, each in {none|floppy|dvd|disk|net}
}

// Refresh reloads the machine information.
func (m *machine) Refresh() error {
	id := m.name
	if id == "" {
		id = m.uUID
	}
	mm, err := getMachine(id)
	if err != nil {
		return err
	}
	*m = *mm
	return nil
}

// Start starts the machine.
func (m *machine) Start() error {
	switch m.state {
	case Paused:
		return vbm("controlvm", m.name, "resume")
	case Poweroff, Saved, Aborted:
		return vbm("startvm", m.name, "--type", "headless")
	}
	return nil
}

// Suspend suspends the machine and saves its state to disk.
func (m *machine) Save() error {
	switch m.state {
	case Paused:
		if err := m.Start(); err != nil {
			return err
		}
	case Poweroff, Aborted, Saved:
		return nil
	}
	return vbm("controlvm", m.name, "savestate")
}

// Pause pauses the execution of the machine.
func (m *machine) Pause() error {
	switch m.state {
	case Paused, Poweroff, Aborted, Saved:
		return nil
	}
	return vbm("controlvm", m.name, "pause")
}

// Stop gracefully stops the machine.
func (m *machine) Stop() error {
	switch m.state {
	case Poweroff, Aborted, Saved:
		return nil
	case Paused:
		if err := m.Start(); err != nil {
			return err
		}
	}

	for m.state != Poweroff { // busy wait until the machine is stopped
		if err := vbm("controlvm", m.name, "acpipowerbutton"); err != nil {
			return err
		}
		time.Sleep(1 * time.Second)
		if err := m.Refresh(); err != nil {
			return err
		}
	}
	return nil
}

// Poweroff forcefully stops the machine. State is lost and might corrupt the disk image.
func (m *machine) Poweroff() error {
	switch m.state {
	case Poweroff, Aborted, Saved:
		return nil
	}
	return vbm("controlvm", m.name, "poweroff")
}

// Restart gracefully restarts the machine.
func (m *machine) Restart() error {
	switch m.state {
	case Paused, Saved:
		if err := m.Start(); err != nil {
			return err
		}
	}
	if err := m.Stop(); err != nil {
		return err
	}
	return m.Start()
}

// Reset forcefully restarts the machine. State is lost and might corrupt the disk image.
func (m *machine) Reset() error {
	switch m.state {
	case Paused, Saved:
		if err := m.Start(); err != nil {
			return err
		}
	}
	return vbm("controlvm", m.name, "reset")
}

// Delete deletes the machine and associated disk images.
func (m *machine) Delete() error {
	if err := m.Poweroff(); err != nil {
		return err
	}
	return vbm("unregistervm", m.name, "--delete")
}

// GetMachine finds a machine by its name or UUID.
func getMachine(id string) (*machine, error) {
	stdout, stderr, err := vbmOutErr("showvminfo", id, "--machinereadable")
	if err != nil {
		if reMachineNotFound.FindString(stderr) != "" {
			return nil, ErrMachineNotExist
		}
		return nil, err
	}
	s := bufio.NewScanner(strings.NewReader(stdout))
	m := &machine{}
	for s.Scan() {
		res := reVMInfoLine.FindStringSubmatch(s.Text())
		if res == nil {
			continue
		}
		key := res[1]
		if key == "" {
			key = res[2]
		}
		val := res[3]
		if val == "" {
			val = res[4]
		}

		switch key {
		case "name":
			m.name = val
		case "UUID":
			m.uUID = val
		case "VMState":
			m.state = MachineState(val)
		case "memory":
			n, err := strconv.ParseUint(val, 10, 32)
			if err != nil {
				return nil, err
			}
			m.memory = uint(n)
		case "cpus":
			n, err := strconv.ParseUint(val, 10, 32)
			if err != nil {
				return nil, err
			}
			m.cPUs = uint(n)
		case "vram":
			n, err := strconv.ParseUint(val, 10, 32)
			if err != nil {
				return nil, err
			}
			m.vRAM = uint(n)
		case "CfgFile":
			m.cfgFile = val
			m.baseFolder = filepath.Dir(val)
		}
	}
	if err := s.Err(); err != nil {
		return nil, err
	}
	return m, nil
}

// ListMachines lists all registered machines.
func listMachines() ([]*machine, error) {
	out, err := vbmOut("list", "vms")
	if err != nil {
		return nil, err
	}
	ms := []*machine{}
	s := bufio.NewScanner(strings.NewReader(out))
	for s.Scan() {
		res := reVMNameUUID.FindStringSubmatch(s.Text())
		if res == nil {
			continue
		}
		m, err := getMachine(res[1])
		if err != nil {
			return nil, err
		}
		ms = append(ms, m)
	}
	if err := s.Err(); err != nil {
		return nil, err
	}
	return ms, nil
}

// CreateMachine creates a new machine. If basefolder is empty, use default.
func createMachine(name, basefolder string) (*machine, error) {
	if name == "" {
		return nil, fmt.Errorf("machine name is empty")
	}

	// Check if a machine with the given name already exists.
	ms, err := listMachines()
	if err != nil {
		return nil, err
	}
	for _, m := range ms {
		if m.name == name {
			return nil, ErrMachineExist
		}
	}

	// Create and register the machine.
	args := []string{"createvm", "--name", name, "--register"}
	if basefolder != "" {
		args = append(args, "--basefolder", basefolder)
	}
	if err := vbm(args...); err != nil {
		return nil, err
	}

	m, err := getMachine(name)
	if err != nil {
		return nil, err
	}

	return m, nil
}

// Modify changes the settings of the machine.
func (m *machine) Modify() error {
	args := []string{"modifyvm", m.name,
		"--firmware", "bios",
		"--bioslogofadein", "off",
		"--bioslogofadeout", "off",
		"--bioslogodisplaytime", "0",
		"--biosbootmenu", "disabled",

		"--ostype", m.oSType,
		"--cpus", fmt.Sprintf("%d", m.cPUs),
		"--memory", fmt.Sprintf("%d", m.memory),
		"--vram", fmt.Sprintf("%d", m.vRAM),

		"--acpi", m.flag.Get(F_acpi),
		"--ioapic", m.flag.Get(F_ioapic),
		"--rtcuseutc", m.flag.Get(F_rtcuseutc),
		"--cpuhotplug", m.flag.Get(F_cpuhotplug),
		"--pae", m.flag.Get(F_pae),
		"--longmode", m.flag.Get(F_longmode),
		"--synthcpu", m.flag.Get(F_synthcpu),
		"--hpet", m.flag.Get(F_hpet),
		"--hwvirtex", m.flag.Get(F_hwvirtex),
		"--triplefaultreset", m.flag.Get(F_triplefaultreset),
		"--nestedpaging", m.flag.Get(F_nestedpaging),
		"--largepages", m.flag.Get(F_largepages),
		"--vtxvpid", m.flag.Get(F_vtxvpid),
		"--vtxux", m.flag.Get(F_vtxux),
		"--accelerate3d", m.flag.Get(F_accelerate3d),
	}

	for i, dev := range m.bootOrder {
		if i > 3 {
			break // Only four slots `--boot{1,2,3,4}`. Ignore the rest.
		}
		args = append(args, fmt.Sprintf("--boot%d", i+1), dev)
	}
	if err := vbm(args...); err != nil {
		return err
	}
	return m.Refresh()
}

// AddNATPF adds a NAT port forarding rule to the n-th NIC with the given name.
func (m *machine) AddNATPF(n int, name string, rule PFRule) error {
	return vbm("controlvm", m.name, fmt.Sprintf("natpf%d", n),
		fmt.Sprintf("%s,%s", name, rule.Format()))
}

// DelNATPF deletes the NAT port forwarding rule with the given name from the n-th NIC.
func (m *machine) DelNATPF(n int, name string) error {
	return vbm("controlvm", m.name, fmt.Sprintf("natpf%d", n), "delete", name)
}

// SetNIC set the n-th NIC.
func (m *machine) SetNIC(n int, nic NIC) error {
	args := []string{"modifyvm", m.name,
		fmt.Sprintf("--nic%d", n), string(nic.Network),
		fmt.Sprintf("--nictype%d", n), string(nic.Hardware),
		fmt.Sprintf("--cableconnected%d", n), "on",
	}

	if nic.Network == "hostonly" {
		args = append(args, fmt.Sprintf("--hostonlyadapter%d", n), nic.HostonlyAdapter)
	}
	return vbm(args...)
}

// AddStorageCtl adds a storage controller with the given name.
func (m *machine) AddStorageCtl(name string, ctl StorageController) error {
	args := []string{"storagectl", m.name, "--name", name}
	if ctl.SysBus != "" {
		args = append(args, "--add", string(ctl.SysBus))
	}
	if ctl.Ports > 0 {
		args = append(args, "--portcount", fmt.Sprintf("%d", ctl.Ports))
	}
	if ctl.Chipset != "" {
		args = append(args, "--controller", string(ctl.Chipset))
	}
	args = append(args, "--hostiocache", bool2string(ctl.HostIOCache))
	args = append(args, "--bootable", bool2string(ctl.Bootable))
	return vbm(args...)
}

// DelStorageCtl deletes the storage controller with the given name.
func (m *machine) DelStorageCtl(name string) error {
	return vbm("storagectl", m.name, "--name", name, "--remove")
}

// AttachStorage attaches a storage medium to the named storage controller.
func (m *machine) AttachStorage(ctlName string, medium StorageMedium) error {
	return vbm("storageattach", m.name, "--storagectl", ctlName,
		"--port", fmt.Sprintf("%d", medium.Port),
		"--device", fmt.Sprintf("%d", medium.Device),
		"--type", string(medium.DriveType),
		"--medium", medium.Medium,
	)
}

func (m *machine) Name() string {
	return m.name
}

func (m *machine) UUID() string {
	return m.uUID
}

func (m *machine) State() MachineState {
	return m.state
}

func (m *machine) CPUs() uint {
	return m.cPUs
}

func (m *machine) Memory() uint {
	return m.memory
}

func (m *machine) VRAM() uint {
	return m.vRAM
}

func (m *machine) CfgFile() string {
	return m.cfgFile
}

func (m *machine) BaseFolder() string {
	return m.baseFolder
}

func (m *machine) OSType() string {
	return m.oSType
}

func (m *machine) Flag() Flag {
	return m.flag
}

func (m *machine) BootOrder() []string {
	return m.bootOrder
}

func (m *machine) SetName(name string) {
	m.name = name
}

func (m *machine) SetUUID(uuid string) {
	// TODO add validation
	m.uUID = uuid
}

func (m *machine) SetState(state MachineState) {
	m.state = state
}

func (m *machine) SetCPUs(cpus uint) {
	m.cPUs = cpus
}

func (m *machine) SetMemory(memory uint) {
	m.memory = memory
}

func (m *machine) SetVRAM(vram uint) {
	m.vRAM = vram
}

func (m *machine) SetCfgFile(cfgFile string) {
	m.cfgFile = cfgFile
}

func (m *machine) SetBaseFolder(baseFolder string) {
	m.baseFolder = baseFolder
}

func (m *machine) SetOSType(osType string) {
	m.oSType = osType
}

func (m *machine) SetFlag(flag Flag) {
	m.flag = flag
}

func (m *machine) SetBootOrder(bootOrder []string) {
	m.bootOrder = bootOrder
}
