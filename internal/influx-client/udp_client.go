package influx_client

import (
	influx_util "github.com/ashirko/influxdb-write-test/internal/influx-util"
	client "github.com/influxdata/influxdb1-client/v2"
	"log"
	"sync"
	"time"
)

func StartClient(wg *sync.WaitGroup, addr string, mesNum, id int) {
	defer wg.Done()
	config := client.UDPConfig{Addr: addr}
	c, err := client.NewUDPClient(config)
	if err != nil {
		log.Println("Error: ", err.Error())
		return
	}
	defer influx_util.CloseAndLog(c)
	bpconfig := client.BatchPointsConfig{Precision: "ms"}
	for i := 0; i < mesNum; i++ {
		data := formData(id, i)
		bp, err := client.NewBatchPoints(bpconfig)
		if err != nil {
			log.Println("Error:", err.Error())
			return
		}
		bp.AddPoint(data)
		if err = c.Write(bp); err != nil {
			log.Println("Error:", err.Error())
		}
		time.Sleep(1 * time.Second)
	}
}
