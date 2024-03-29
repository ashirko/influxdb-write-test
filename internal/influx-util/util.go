package influx_util

import (
	client "github.com/influxdata/influxdb1-client/v2"
	"io"
	"log"
	"sync"
	"time"
)

func Nanoseconds() int64 {
	return time.Now().UnixNano()
}

func CloseAndLog(c io.Closer) {
	err := c.Close()
	if err != nil {
		log.Println("Error:", err.Error())
	}
}

func SentToBuff(ch chan *client.Point, m *client.Point) {
	select {
	case ch <- m:
	default:
		log.Println("Error: buffer is full!")
	}
}

func UdpClient(addr string) client.Client {
	config := client.UDPConfig{Addr: addr}
	c, err := client.NewUDPClient(config)
	if err != nil {
		log.Println("Error: ", err.Error())
		return nil
	}
	return c
}

func HttpClient(addr string) client.Client {
	config := client.HTTPConfig{Addr: addr}
	c, err := client.NewHTTPClient(config)
	if err != nil {
		log.Println("Error: ", err.Error())
		return nil
	}
	return c
}

func BpConfig(db string) client.BatchPointsConfig {
	bpconfig := client.BatchPointsConfig{
		Precision: "ms",
		Database:  db,
	}
	return bpconfig
}

func StartSender(c client.Client, wgBuff *sync.WaitGroup, buff chan *client.Point, ch chan bool, sentPeriod, maxPoints int, bpconfig client.BatchPointsConfig) {
	defer wgBuff.Done()
	bp, err := client.NewBatchPoints(bpconfig)
	ticker := time.NewTicker(time.Duration(sentPeriod) * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ch:
			if len(bp.Points()) > 0 {
				SentData(c, bp)
			}
			return
		case <-ticker.C:
			if len(bp.Points()) > 0 {
				SentData(c, bp)
				bp, err = client.NewBatchPoints(bpconfig)
				if err != nil {
					log.Println("Error:", err.Error())
					return
				}
			}
		case Point := <-buff:
			bp.AddPoint(Point)
			if len(bp.Points()) == maxPoints {
				SentData(c, bp)
				bp, err = client.NewBatchPoints(bpconfig)
				if err != nil {
					log.Println("Error:", err.Error())
					return
				}
			}
		}
	}
}

func SentData(c client.Client, bp client.BatchPoints) {
	if err := c.Write(bp); err != nil {
		log.Println("Error:", err.Error())
	}
}
