package main

import (
	"math/rand"
	"time"

	. "github.com/PMoneda/telemetry"
	"github.com/PMoneda/telemetry/registry"
)

func main() {
	config := registry.Config{}
	config.Host = "http://localhost"
	config.Port = "8086"
	value := Database("test").RetentionPolicy("autogen").Measurement("sensor").Tag("floor").Value("1st")
	t := BuildTelemetryContext(config, Context(value))

	rand.Seed(int64(time.Now().Nanosecond()))
	runtime := Database("test").RetentionPolicy("autogen").Measurement("metrics").Tag("host").Value("host0")
	telemetryRuntime := BuildTelemetryContext(config, Context(runtime))
	go telemetryRuntime.StartRuntimeTelemetry()
	go t.StartTelemetry()
	for {
		t.Push("success", rand.Float64()*500)
		t.Push("failure", rand.Float64()*10)
		time.Sleep(1 * time.Second)
	}
}
