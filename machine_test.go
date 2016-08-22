package virtualbox

import "testing"

// compile time check to make sure machine type implements Machine interface
var _  Machine = (*machine)(nil)

func TestMachine(t *testing.T) {
	ms, err := listMachines()
	if err != nil {
		t.Fatal(err)
	}
	for _, m := range ms {
		t.Logf("%+v", m)
	}
}
