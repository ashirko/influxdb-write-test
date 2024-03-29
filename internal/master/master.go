package master

import (
	"flag"
	"github.com/ashirko/influxdb-write-test/internal/influx-client"
	influx_util "github.com/ashirko/influxdb-write-test/internal/influx-util"
	client "github.com/influxdata/influxdb1-client/v2"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type ScriptParams struct {
	ConNum, MesNum int
	timeWait int
	ScriptName     string
	UdpAddress     string
	HttpAddress    string
}

const conNumDef = 100
const mesNumDef = 100
const udpAddrDef = "127.0.0.1:8089"
const httpAddrDef = "http://localhost:8086"
const waitDefault = 5
const UdpSentPeriod = 5
const UdpMaxPoints = 10
const HttpSentPeriod = 5
const HttpMaxPoints = 100
const UdpDB = "udp"
const HttpDB = "http"

var DBName string

func StartTest() {
	params := parseFlags()
	log.Printf("start %s for %d connections", params.ScriptName, params.ConNum)
	tStart := influx_util.Nanoseconds()/1000000*1000000
	if params.ScriptName == "influx-test" {
		DBName = UdpDB
		runUdp(params)
	} else if params.ScriptName == "influx-test-buff" {
		DBName = UdpDB
		runUdpBuff(params)
	} else if params.ScriptName == "influx-test-http" {
		DBName = HttpDB
		runHttpBuff(params)
	} else {
		log.Println("Error: test function doesn't exits")
		os.Exit(1)
	}
	time.Sleep(time.Duration(params.timeWait) * time.Second)
	tFinish := influx_util.Nanoseconds()
	data := requestData(params, tStart, tFinish)
	compareResult(params, data)
}

func runUdp(params *ScriptParams) {
	var wg sync.WaitGroup
	startUdpClients(params, &wg)
	wg.Wait()
}

func startUdpClients(params *ScriptParams, wg *sync.WaitGroup) {
	for i := 1; i <= params.ConNum; i++ {
		wg.Add(1)
		go influx_client.StartClient(wg, params.UdpAddress, params.MesNum, i)
		time.Sleep(1 * time.Millisecond)
	}
}

func runUdpBuff(params *ScriptParams) {
	var wg, wgBuff sync.WaitGroup
	buff := make(chan *client.Point, 20000)
	ch := make(chan bool)
	c := influx_util.UdpClient(params.UdpAddress)
	defer influx_util.CloseAndLog(c)
	wgBuff.Add(1)
	bpconfig := influx_util.BpConfig(DBName)
	go influx_util.StartSender(c, &wgBuff, buff, ch, UdpSentPeriod, UdpMaxPoints, bpconfig)
	startClientsBuff(params, &wg, buff)
	wg.Wait()
	close(ch)
	wgBuff.Wait()
}

func runHttpBuff(params *ScriptParams) {
	var wg, wgBuff sync.WaitGroup
	buff := make(chan *client.Point, 20000)
	ch := make(chan bool)
	c := influx_util.HttpClient(params.HttpAddress)
	defer influx_util.CloseAndLog(c)
	wgBuff.Add(1)
	bpconfig := influx_util.BpConfig(DBName)
	go influx_util.StartSender(c, &wgBuff, buff, ch, HttpSentPeriod, HttpMaxPoints, bpconfig)
	startClientsBuff(params, &wg, buff)
	wg.Wait()
	close(ch)
	wgBuff.Wait()
}


func startClientsBuff(params *ScriptParams, wg *sync.WaitGroup, buff chan *client.Point) {
	for i := 1; i <= params.ConNum; i++ {
		wg.Add(1)
		go influx_client.StartClientBuff(buff, wg, params.MesNum, i)
		time.Sleep(1 * time.Millisecond)
	}
}
func parseFlags() *ScriptParams {
	vals := new(ScriptParams)
	flag.IntVar(&vals.ConNum, "c", conNumDef, "number of connections")
	flag.IntVar(&vals.MesNum, "m", mesNumDef, "number of messages")
	flag.IntVar(&vals.timeWait, "s", waitDefault, "wait before checking result")
	flag.StringVar(&vals.UdpAddress, "u", udpAddrDef, "InfluxDB UDP Address")
	flag.StringVar(&vals.HttpAddress, "h", httpAddrDef, "InfluxDB HTTP Address")
	flag.Parse()
	vals.ScriptName = filepath.Base(os.Args[0])
	return vals
}
