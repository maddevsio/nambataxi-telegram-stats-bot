package main

import (
	"time"
	"encoding/json"
	sc "github.com/maddevsio/simple-config"
	"gopkg.in/resty.v1"
	"log"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestRequest(t *testing.T) {
	config := sc.NewSimpleConfig("./config", "yml")
	url := config.GetString("url")
	resp, err := resty.R().Get(url)

	checkErr(err)

	log.Print(url)
	log.Print(resp)
}

type logData struct {
	Target     string      `json:"target"`
	Datapoints [][]float64 `json:"datapoints"`
}

func GetMaxDataFromJSON(raw string) int {
	var data []logData
	_ = json.Unmarshal([]byte(raw), &data)

	var picked []int
	var max int
	for _, v := range data[0].Datapoints {
		if v[0] > 0 {
			if int(v[0]) > max {
				max = int(v[0])
			}
			picked = append(picked, int(v[0]))
		}
	}
	return max
}

func TestParseResult(t *testing.T) {
	config := sc.NewSimpleConfig("./config", "yml")
	url := config.GetString("url")
	resp, err := resty.R().Get(url)

	checkErr(err)

	log.Print(GetMaxDataFromJSON(resp.String()))
}

func TestGetDataDayBefore(t *testing.T) {
	timeForTest := time.Date(2017,11,1,0,0,0,0,time.UTC)
	dayBefore := GetDayBeforeInFormat(timeForTest)
	assert.Equal(t, "20171031", dayBefore)
}

func checkErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}
