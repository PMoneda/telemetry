package sensors

import "github.com/PMoneda/telemetry/registry"

//Sensor is a basic interface to represent a metric sensor
type Sensor interface {
	Plug(registry.Registry, string)
	Read()
}
