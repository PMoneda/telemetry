package main

import (
	"github.com/PMoneda/telemetry"
	"github.com/PMoneda/telemetry/registry"
)

func main() {
	config := registry.Config{}
	config.Host = "http://localhost"
	config.Port = "8086"
	value := telemetry.Database("my-metrics").RetentionPolicy("autogen").Measurement("network-access").Tag("url").Value("http://pruu.herokuapp.com/dump/test-network-beat")
	telemetry.BuildTelemetryContext(config, telemetry.Context(value))
	//t.Listen()
}
