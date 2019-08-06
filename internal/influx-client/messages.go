package influx_client

import (
	client "github.com/influxdata/influxdb1-client/v2"
	"log"
	"strconv"
	"time"
)

const measurement = "monitoring1"

func formData(id, mesNum int) *client.Point {
	tags := map[string]string{
		"id":          strconv.Itoa(id),
		"ip":          "11.192.12." + strconv.Itoa(id),
		"service_id":  "2",
		"packet_type": "1",
		"direction":   "from_bnst",
	}
	fields := map[string]interface{}{
		"session_id": "abcdef",
		"bytes":      216,
		"npl_id":     mesNum,
		"nph_id":     mesNum,
		"packet_id":  mesNum + 1,
		"dump":       false,
		"double":     false,
	}
	pt, err := client.NewPoint(measurement, tags, fields, time.Now())
	if err != nil {
		log.Println("Error: ", err.Error())
	}
	return pt
}
