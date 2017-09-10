package sensors

import "github.com/PMoneda/telemetry/registry"

//BaseSensor is a base struct for all sensors
type BaseSensor struct {
	context  string
	registry registry.Registry
}

//Plug a registry and a name context to this sensor
func (sensor *BaseSensor) Plug(reg registry.Registry, context string) {
	sensor.context = context
	sensor.registry = reg
}
func (sensor *BaseSensor) register(tag string, value interface{}) {
	if sensor.context != "" {
		sensor.registry.Register(sensor.context+":::"+tag, value)
	} else {
		sensor.registry.Register(tag, value)
	}
}

func (sensor *BaseSensor) Read() {
	//empty implementation for sensors without passive reads
}
