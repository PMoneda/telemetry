package sensors

import (
	"crypto/tls"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/PMoneda/telemetry/registry"
)

var defaultDialer = &net.Dialer{Timeout: 16 * time.Second, KeepAlive: 16 * time.Second}

var cfg *tls.Config = &tls.Config{
	InsecureSkipVerify: true,
}
var client *http.Client = &http.Client{
	Transport: &http.Transport{
		TLSClientConfig:     cfg,
		Dial:                defaultDialer.Dial,
		TLSHandshakeTimeout: 16 * time.Second,
	},
}

func doRequest(method, url, body string, header map[string]string) (string, int, error) {
	message := strings.NewReader(body)
	req, err := http.NewRequest(method, url, message)
	if err != nil {
		return "", http.StatusInternalServerError, err
	}
	if header != nil {
		for k, v := range header {
			req.Header.Add(k, v)
		}
	}
	resp, errResp := client.Do(req)
	if errResp != nil {
		return "", 0, errResp
	}
	defer resp.Body.Close()
	data, errResponse := ioutil.ReadAll(resp.Body)
	if errResponse != nil {
		return "", resp.StatusCode, errResponse
	}
	sData := string(data)
	return sData, resp.StatusCode, nil
}

//NetworkBeatSensor is a sensor to access some URL and collect metrics
type NetworkBeatSensor struct {
	BaseSensor
	header map[string]string
	body   string
	url    string
	method string
}

//Plug a registry and a context to sensor
func (sensor *NetworkBeatSensor) Plug(reg registry.Registry, context string) {
	sensor.context = context
	sensor.registry = reg
}

//Config inputs to sensor
func (sensor *NetworkBeatSensor) Config(url string, method string, body string, header map[string]string) {
	sensor.url = url
	sensor.body = body
	sensor.header = header
	sensor.method = method
}

//Read data from sensor and writes to registry
func (sensor *NetworkBeatSensor) Read() {
	start := time.Now()
	response, status, err := doRequest(sensor.method, sensor.url, sensor.body, sensor.header)
	end := time.Now()
	duration := end.Sub(start)
	sensor.register("response", response)
	sensor.register("status", status)
	if err != nil {
		sensor.register("error", err.Error())
	}
	sensor.register("request-time", duration.Seconds())
}
