package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"gopkg.in/resty.v1"
	"testing"
	"time"
)

func TestSendTelegramMessage(t *testing.T) {
	config.Fill("./config", "yaml")
	message := "Test"
	err := ConnectTelegramAndSendMessage(message, config)
	assert.NoError(t, err)
	config.Token = "non existant Telegram token"
	err = ConnectTelegramAndSendMessage(message, config)
	assert.Error(t, err)
}

func TestGetMaxForDateAndTarget(t *testing.T) {
	config.Fill("./config", "yaml")
	assert.Equal(t, "9583", GetMaxForDateAndTarget("20171031", "taxi.orders.total", config))
}

func TestRequest(t *testing.T) {
	config.Fill("./config", "yaml")
	resp, err := resty.R().Get(config.Url)
	checkErr(err)
	assert.Equal(t, "200 OK", resp.Status())
}

func TestParseResult(t *testing.T) {
	config.Fill("./config", "yaml")
	resp, err := resty.R().Get(fmt.Sprintf(config.Url, "20171031", "taxi.orders.total"))
	checkErr(err)
	assert.Equal(t, 9583, GetMaxDataFromJSON(resp.String()))
}

func TestGetDataDayBefore(t *testing.T) {
	timeForTest := time.Date(2017, 11, 1, 0, 0, 0, 0, time.UTC)
	dayBefore := GetDayBeforeInFormat(timeForTest)
	assert.Equal(t, "20171031", dayBefore)
}
