package influxdb

import (
	"strings"
	"sync"
	"time"

	client "github.com/influxdata/influxdb/client/v2"
)

const sep string = ":::"

//InfluxDB is the interface to Influx instance
type InfluxDB struct {
	Port      string
	Host      string
	Precision string
	Client    client.Client
	Buffer    map[string]client.BatchPoints
	mutex     sync.Mutex
}

type registryMetaData struct {
	DB          string
	Retention   string
	Measurement string
	Tag         string
	TagValue    string
	Field       string
	FieldValue  interface{}
}

//New returns a new influxdb connection
func New(host, port string) *InfluxDB {
	db := new(InfluxDB)
	db.Host = host
	db.Port = port
	db.Precision = "s"
	db.Buffer = make(map[string]client.BatchPoints)
	db.Connect()
	return db
}

//Register a new Point to InfluxDB
func (db *InfluxDB) Register(tag string, value interface{}) (err error) {
	// Create a new point batch
	meta := db.parseTag(tag)
	bp, err := db.FindOrCreateBatchPoint(tag)
	if err != nil {
		return
	}
	tags := map[string]string{meta.Tag: meta.TagValue}
	fields := map[string]interface{}{
		meta.Field: value,
	}

	pt, err := client.NewPoint(meta.Measurement, tags, fields, time.Now())
	bp.AddPoint(pt)
	if err != nil {
		return
	}
	return
}

//Flush save data to influxdb
func (db *InfluxDB) Flush(tag string) (err error) {
	meta := db.parseTag(tag)
	bp, found := db.Buffer[meta.DB+sep+meta.Retention]
	if found {
		err = db.Client.Write(bp)
	}
	return
}

//FindOrCreateBatchPoint creates a new bach points if not exist
func (db *InfluxDB) FindOrCreateBatchPoint(tag string) (batch client.BatchPoints, err error) {
	meta := db.parseTag(tag)
	batch, found := db.Buffer[meta.DB+sep+meta.Retention]
	if !found {
		batch, err = db.createBatchPoint(tag)
	}
	return
}

func (db *InfluxDB) createBatchPoint(tag string) (batch client.BatchPoints, err error) {
	db.mutex.Lock()
	defer db.mutex.Unlock()
	meta := db.parseTag(tag)
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:        meta.DB,
		Precision:       db.Precision,
		RetentionPolicy: meta.Retention,
	})
	db.Buffer[meta.DB+sep+meta.Retention] = bp
	batch = bp
	return
}

func (db *InfluxDB) parseTag(tag string) (meta registryMetaData) {
	parts := strings.Split(tag, sep)
	meta.DB = parts[0]
	meta.Retention = parts[1]
	meta.Measurement = parts[2]
	meta.Tag = parts[3]
	meta.TagValue = parts[4]
	meta.Field = parts[5]
	return
}

//Connect stablished connection with influxdb
func (db *InfluxDB) Connect() (err error) {
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: db.Host + ":" + db.Port,
	})
	if err != nil {
		return
	}
	db.Client = c
	return
}

//FlushAll save all batche points to influxdb
func (db *InfluxDB) FlushAll() (err error) {
	db.mutex.Lock()
	defer db.mutex.Unlock()
	for k, v := range db.Buffer {
		err = db.Client.Write(v)
		if err != nil {
			return
		}
		delete(db.Buffer, k)
	}
	return
}

//IsConnected checks if connection is open
func (db *InfluxDB) IsConnected() (connected bool) {
	_, _, err := db.Client.Ping(200 * time.Millisecond)
	if err != nil {
		connected = false
	}
	connected = true
	return
}
