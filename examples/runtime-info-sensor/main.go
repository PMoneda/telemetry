package main

import (
	"time"

	"github.com/PMoneda/telemetry"
	"github.com/PMoneda/telemetry/registry"
	"github.com/PMoneda/telemetry/sensors"
)

func main() {
	config := registry.Config{}
	config.Host = "http://localhost"
	config.Port = "8086"
	value := telemetry.Database("my-metrics").RetentionPolicy("autogen").Measurement("runtime").Tag("host").Value("host0")
	tel := telemetry.BuildTelemetryContext(config, telemetry.Context(value))

	pruu := new(sensors.RuntimeSensor)
	tel.Listen(pruu, 1*time.Second)
	w := make(chan int)
	<-w
}
