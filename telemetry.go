package telemetry

import (
	"time"

	"github.com/PMoneda/telemetry/registry"
	"github.com/PMoneda/telemetry/sensors"
)

//Telemetry is a struct to manage sensor and collect data
type Telemetry struct {
	Registry registry.Registry
	Root     string
}

var sep string = ":::"

//Context is a basic type to config telemetry
type Context func(*Telemetry) *Telemetry

//RetentionPolicy to config retention policy on influxdb
type RetentionPolicy Context

//Measurement config on influxdb
type Measurement Context

//Tag name on influxdb
type Tag Context

//Value is a tag value on influxdb
type Value Context

//Tag set tag name
func (r Tag) Tag(value string) Value {
	return func(t *Telemetry) *Telemetry {
		r(t)
		t.Root = t.Root + sep + value
		return t
	}
}

//Value set value
func (r Value) Value(value string) Context {
	return func(t *Telemetry) *Telemetry {
		r(t)
		t.Root = t.Root + sep + value
		return t
	}
}

//Measurement set measurement name
func (r Measurement) Measurement(s string) Tag {
	return func(t *Telemetry) *Telemetry {
		r(t)
		t.Root = t.Root + sep + s
		return t
	}
}

//RetentionPolicy set name
func (r RetentionPolicy) RetentionPolicy(s string) Measurement {
	return func(t *Telemetry) *Telemetry {
		r(t)
		t.Root = t.Root + sep + s
		return t
	}
}

//Database set name
func Database(s string) RetentionPolicy {
	return (func(t *Telemetry) *Telemetry {
		t.Root = s
		return t
	})
}

//NewTelemetryForInfluxDB creates a new telemetry to influxdb
func NewTelemetryForInfluxDB(config registry.Config, context string) *Telemetry {
	telemetry := new(Telemetry)
	telemetry.Root = context
	telemetry.Registry = registry.NewInfluxClient(config)
	return telemetry
}

//BuildTelemetryContext build a new telemetry based on context
func BuildTelemetryContext(config registry.Config, ctx Context) *Telemetry {
	telemetry := new(Telemetry)
	ctx(telemetry)
	telemetry.Registry = registry.NewInfluxClient(config)
	return telemetry
}

//Push data to this telemetry
func (t *Telemetry) Push(tag string, value interface{}) (err error) {
	err = t.Registry.Register(t.Root+sep+tag, value)
	return
}

//Flush flushes all registry
func (t *Telemetry) Flush() (err error) {
	err = t.Registry.FlushAll()
	return
}

//Listen to a specific sensor and flush sensor reads
func (t *Telemetry) Listen(sensor sensors.Sensor, duration time.Duration) {
	tick := time.Tick(duration)
	for {
		select {
		case <-tick:
			sensor.Plug(t.Registry, t.Root)
			sensor.Read()
			t.Flush()
		}
	}
}

func (t *Telemetry) StartTelemetry() {
	tick := time.Tick(1 * time.Second)
	for {
		select {
		case <-tick:
			t.Flush()
		}
	}
}
