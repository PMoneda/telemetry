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
	value := telemetry.Database("my-metrics").RetentionPolicy("autogen").Measurement("network-access").Tag("url").Value("http://pruu.herokuapp.com/dump/test-network-beat")
	tel := telemetry.BuildTelemetryContext(config, telemetry.Context(value))

	pruu := new(sensors.NetworkBeatSensor)
	pruu.Config("https://pruu.herokuapp.com/dump/test-telemetry-pruu", "POST", "test", nil)
	tel.Listen(pruu, 1*time.Second)
	w := make(chan int)
	<-w
}
