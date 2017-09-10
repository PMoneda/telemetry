package sensors

import "testing"

type mockRegistry struct {
	*testing.T
}

func (m *mockRegistry) Register(tag string, value interface{}) error {
	switch tag {
	case "memory-alloc", "total-memory-alloc", "garbage-collector-total-run", "num-goroutines", "num-cgo-calls":
	default:
		m.Fail()
	}
	return nil
}
func (m *mockRegistry) Flush(tag string) error {
	return nil
}

func (m *mockRegistry) FlushAll() error {
	return nil
}

func TestRuntimeSensor(t *testing.T) {
	rtSensor := RuntimeSensor{}
	reg := new(mockRegistry)
	reg.T = t
	rtSensor.Plug(reg, "")
	rtSensor.Read()

}
