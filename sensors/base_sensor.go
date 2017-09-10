package sensors

import "github.com/PMoneda/telemetry/registry"

//BaseSensor is a base struct for all sensors
type BaseSensor struct {
	context  string
	registry registry.Registry
}

func (sensor *BaseSensor) register(tag string, value interface{}) {
	if sensor.context != "" {
		sensor.registry.Register(sensor.context+":::"+tag, value)
	} else {
		sensor.registry.Register(tag, value)
	}

}
