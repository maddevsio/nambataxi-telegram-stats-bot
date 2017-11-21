package main

import (
	"testing"
	"log"
	"gopkg.in/resty.v1"
	"encoding/json"
	sc "github.com/maddevsio/simple-config"
)

func TestRequest(t *testing.T) {
	config    := sc.NewSimpleConfig("./config", "yml")
	url       := config.GetString("url")
	resp, err := resty.R().Get(url)

	checkErr(err)

	log.Print(url)
	log.Print(resp)
}

type logData struct {
	Target string `json:"target"`
	Datapoints [][]float64 `json:"datapoints"`
}

func TestParseResult(t *testing.T) {
	config    := sc.NewSimpleConfig("./config", "yml")
	url       := config.GetString("url")
	resp, err := resty.R().Get(url)

	checkErr(err)

	var data []logData
	_ = json.Unmarshal([]byte(resp.String()), &data)

	var picked []int
	var max int
	for _, v := range data[0].Datapoints {
		log.Print(v[0])
		if v[0] > 0 {
			if int(v[0]) > max {
				max = int(v[0])
			}
			picked = append(picked, int(v[0]))
		}
	}
	log.Print(picked)
	log.Print(max)
}

func checkErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}
