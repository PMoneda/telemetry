package sensors

import (
	"time"
)

//TimingSensor measures total time of execution
type TimingSensor struct {
	BaseSensor
}

func (sensor *TimingSensor) Attach(label string, callback func()) {
	start := time.Now()
	callback()
	end := time.Now()
	duration := end.Sub(start)
	sensor.register(label, duration.Seconds())
}
