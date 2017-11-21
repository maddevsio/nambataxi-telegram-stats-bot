package main

import (
	"testing"
	"log"
	"gopkg.in/resty.v1"
	sc "github.com/maddevsio/simple-config"
)

func TestRequest(t *testing.T) {
	config := sc.NewSimpleConfig("./config", "yml")
	url := config.GetString("url")
	log.Print(url)
	resp, err := resty.R().Get(url)
	checkErr(err)
	log.Print(resp)
}

func checkErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}
