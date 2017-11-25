package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"gopkg.in/resty.v1"
	"testing"
	"time"
	"os"
)

func TestGetPicAboutCabs(t *testing.T) {
	path := "/tmp/driversTestGetPicAboutCabs.png"
	_ = os.Remove(path)
	validPngInfo := path + ": PNG image data, 586 x 308, 8-bit/color RGB, non-interlaced\n"
	config.Fill("./config", "yaml")
	date := GetDayBeforeInFormat(time.Now())
	err := GetPicAboutCabs(date, path, config)
	assert.NoError(t, err)
	imgInfo := exe("file", []string{path})
	assert.Equal(t, validPngInfo, imgInfo)
	err = os.Remove(path)
	assert.NoError(t, err)
}

func TestSendPicToTelegramChat(t *testing.T) {
	config.Fill("./config", "yaml")
	path := "/tmp/driversTestSendPicToTelegramChat.png"
	_ = os.Remove(path)
	date := GetDayBeforeInFormat(time.Now())
	err := GetPicAboutCabs(date, path, config)
	assert.NoError(t, err)
	err = ConnectTelegramAndSendPic(path, "Распределение машин за вчера", config)
	assert.NoError(t, err)
}

func TestGetFreeCabsNamba(t *testing.T) {
	config.Fill("./config", "yaml")
	freeCabs := GetFreeCabsNamba(config)
	assert.NotZero(t, freeCabs)
}

func TestSendFullInfo(t *testing.T) {
	config.Fill("./config", "yaml")
	SendFullInfo(config)
}

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

func TestParseResultOrders(t *testing.T) {
	config.Fill("./config", "yaml")
	resp, err := resty.R().Get(fmt.Sprintf(config.Url, "20171031", "taxi.orders.total"))
	checkErr(err)
	assert.Equal(t, 9583, GetMaxDataFromJSON(resp.String()))
}

func TestParseResultOrdersComfort(t *testing.T) {
	config.Fill("./config", "yaml")
	resp, err := resty.R().Get(fmt.Sprintf(config.Url, "20171124", "taxi.orders_comfort.total"))
	checkErr(err)
	assert.Equal(t, 308, GetMaxDataFromJSON(resp.String()))
}

func TestGetDataDayBefore(t *testing.T) {
	timeForTest := time.Date(2017, 11, 1, 0, 0, 0, 0, time.UTC)
	dayBefore := GetDayBeforeInFormat(timeForTest)
	assert.Equal(t, "20171031", dayBefore)
}
