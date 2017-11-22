package main

import (
	"time"
	"gopkg.in/resty.v1"
	"log"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestRequest(t *testing.T) {
	config.Fill("./config", "yaml")
	resp, err := resty.R().Get(config.Url)
	checkErr(err)
	assert.Equal(t, "200 OK", resp.Status())
}

func TestParseResult(t *testing.T) {
	config.Fill("./config", "yaml")
	resp, err := resty.R().Get(config.Url)
	checkErr(err)
	log.Print(GetMaxDataFromJSON(resp.String()))
}

func TestGetDataDayBefore(t *testing.T) {
	timeForTest := time.Date(2017,11,1,0,0,0,0,time.UTC)
	dayBefore := GetDayBeforeInFormat(timeForTest)
	assert.Equal(t, "20171031", dayBefore)
}
