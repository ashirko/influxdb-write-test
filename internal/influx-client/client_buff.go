package influx_client

import (
	influx_util "github.com/ashirko/influxdb-udp-test/internal/influx-util"
	client "github.com/influxdata/influxdb1-client/v2"
	"sync"
	"time"
)

func StartClientBuff(buff chan *client.Point, wg *sync.WaitGroup, mesNum, id int) {
	defer wg.Done()
	for i := 0; i < mesNum; i++ {
		data := formData(id, i)
		influx_util.SentToBuff(buff, data)
		time.Sleep(1 * time.Second)
	}
}
