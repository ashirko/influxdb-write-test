package master

import (
	"encoding/json"
	"fmt"
	influx_util "github.com/ashirko/influxdb-udp-test/internal/influx-util"
	client "github.com/influxdata/influxdb1-client/v2"
	"log"
	"strconv"
)

func requestData(params *ScriptParams, from, to int64) []client.Result {
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: params.HttpAddress,
	})
	if err != nil {
		fmt.Println("Error creating InfluxDB Client: ", err.Error())
		return nil
	}
	defer influx_util.CloseAndLog(c)
	query := "select count(bytes) from monitoring1 where time >= " + strconv.FormatInt(from, 10) + " AND time <= " + strconv.FormatInt(to, 10)
	q := client.NewQuery(query, "testudp", "")
	if response, err := c.Query(q); err == nil && response.Error() == nil {
		return response.Results
	}
	return nil
}

func compareResult(params *ScriptParams, data []client.Result) {
	expected := params.ConNum * params.MesNum
	//log.Println("data:", data)
	res, err := data[0].Series[0].Values[0][1].(json.Number).Int64()
	if err != nil {
		log.Println("Error: ", err.Error())
	} else {
		if res == int64(expected) {
			log.Printf("Success! Wrote %v records\n", res)
		} else {
			log.Printf("Error! Expected %v records, but wrote %v records\n", expected, res)
		}
	}
}
