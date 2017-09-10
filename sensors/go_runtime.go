package sensors

import (
	"runtime"

	"github.com/PMoneda/telemetry/registry"
)

//RuntimeSensor collect data from go runtime
type RuntimeSensor struct {
	BaseSensor
}

//Plug a registry and a name context to this sensor
func (sensor *RuntimeSensor) Plug(reg registry.Registry, context string) {
	sensor.context = context
	sensor.registry = reg
}

func (sensor *RuntimeSensor) Read() {

	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	alloc := float64(m.Alloc) / (1024)
	totalAlloc := float64(m.TotalAlloc) / (1024)
	gcRun := float64(m.NumGC)
	sensor.register("memory-alloc", alloc)
	sensor.register("total-memory-alloc", totalAlloc)
	sensor.register("garbage-collector-total-run", gcRun)
	sensor.register("num-goroutines", runtime.NumGoroutine())
	sensor.register("num-cgo-calls", runtime.NumCgoCall())
}
