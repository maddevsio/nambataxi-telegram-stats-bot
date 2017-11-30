package main

import (
	"github.com/carlescere/scheduler"
	sc "github.com/maddevsio/simple-config"
	"gopkg.in/resty.v1"
	"gopkg.in/telegram-bot-api.v4"

	"io/ioutil"
	"os"
	"os/exec"
	"encoding/json"
	"fmt"
	"log"
	"runtime"
	"strconv"
	"time"
)

type logData struct {
	Target     string      `json:"target"`
	Datapoints [][]float64 `json:"datapoints"`
}

type Config struct {
	Url    string
	Token  string
	ChatID int64
	PicUrl string

	FreeCabsNambaUrl	 string
	AllCabsNambaUrl		 string
	TimeForYesterdayData string
	TimeForEveningCabs	 string
	TimeForDriversData   string
}

type Coord struct {
	Lat string `json:"lat"`
	Lng string `json:"lng"`
}

type Drivers struct {
	Drivers []Coord `json:"drivers"`
}

var config = Config{}

const (
	TARGET_ORDERS_TOTAL    = "taxi.orders.total"
	TARGET_ORDERS_FINISHED = "taxi.orders.finished"
	TARGET_ORDERS_REJECTED = "taxi.orders.rejected"

	TARGET_ORDERS_COMFORT_TOTAL    = "taxi.orders_comfort.total"
	TARGET_ORDERS_COMFORT_FINISHED = "taxi.orders_comfort.finished"
	TARGET_ORDERS_COMFORT_REJECTED = "taxi.orders_comfort.rejected"

	TARGET_DRIVERS_FREE    = "taxi.drivers.free"
	TARGET_DRIVERS_TOTAL   = "taxi.drivers.total"
)

func (cs *Config) Fill(configFile string, configExt string) {
	c := sc.NewSimpleConfig(configFile, configExt)
	cs.Url    = c.GetString("url")
	cs.PicUrl = c.GetString("picurl")
	cs.Token  = c.GetString("token")
	cs.ChatID = int64(c.Get("chatid").(int))
	cs.TimeForYesterdayData = c.GetString("timeforyesterdaydata")
	cs.TimeForEveningCabs	= c.GetString("timeforeveningcabs")
	cs.FreeCabsNambaUrl     = c.GetString("freecabsnambaurl")
	cs.AllCabsNambaUrl      = c.GetString("allcabsnambaurl")
}

func GetDayBeforeInFormat(t time.Time) string {
	return t.AddDate(0, 0, -1).Format("20060102")
}

func GetFreeCabsNamba(config Config) int {
	resp, err := resty.R().Get(config.FreeCabsNambaUrl)
	checkErr(err)
	var drivers Drivers
	err = json.Unmarshal([]byte(resp.String()), &drivers)
	checkErr(err)
	return len(drivers.Drivers)
}

func GetAllCabsNamba(config Config) int {
	resp, err := resty.R().Get(config.AllCabsNambaUrl)
	checkErr(err)
	var drivers Drivers
	err = json.Unmarshal([]byte(resp.String()), &drivers)
	checkErr(err)
	return len(drivers.Drivers)
}

func GetPicAboutCabs(date string, path string, config Config) error {
	url := fmt.Sprintf(config.PicUrl, date, date, TARGET_DRIVERS_FREE, TARGET_DRIVERS_TOTAL)
	resp, err := resty.R().Get(url)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path, resp.Body(), 0644)
	if err != nil {
		return err
	}
	return nil
}

func GetMaxForDateAndTarget(date string, target string, config Config) string {
	url := fmt.Sprintf(config.Url, date, target)
	resp, err := resty.R().Get(url)
	checkErr(err)
	return strconv.Itoa(GetMaxDataFromJSON(resp.String()))
}

func GetMaxDataFromJSON(raw string) int {
	var data []logData
	err := json.Unmarshal([]byte(raw), &data)
	checkErr(err)
	if len(data) == 0 {
		return 0
	}

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

func ConnectTelegramAndSendPic(path string, caption string, config Config) error {
	bot, err := tgbotapi.NewBotAPI(config.Token)
	if err != nil {
		return err
	}
	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self)

	msg := tgbotapi.NewPhotoUpload(config.ChatID, path)
	msg.Caption = caption
	_, err = bot.Send(msg)
	if err != nil {
		return err
	}
	return nil

}

func ConnectTelegramAndSendMessage(message string, config Config) error {
	bot, err := tgbotapi.NewBotAPI(config.Token)
	if err != nil {
		return err
	}
	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	msg   := tgbotapi.NewMessage(config.ChatID, message)
	_, err = bot.Send(msg)
	if err != nil {
		return err
	}
	return nil
}

func GetRejectPercent(maxTotal string, maxRejected string) string {
	maxRejectedInt, err := strconv.Atoi(maxRejected)
	checkErr(err)
	maxTotalInt, err := strconv.Atoi(maxTotal)
	checkErr(err)
	return strconv.Itoa(maxRejectedInt*100/maxTotalInt) + "%"
}

func CreateMessageForOrders(targets [3]string) string {
	var message string

	maxTotal := GetMaxForDateAndTarget(GetDayBeforeInFormat(time.Now()), targets[0], config)
	message += "Всего заказов: " + maxTotal + "\n"

	maxFinished := GetMaxForDateAndTarget(GetDayBeforeInFormat(time.Now()), targets[1], config)
	message += "Всего выполненных заказов: " + maxFinished + "\n"

	maxRejected := GetMaxForDateAndTarget(GetDayBeforeInFormat(time.Now()),	targets[2], config)
	message += "Всего отмененных заказов: " + maxRejected + "\n"

	rejectedPercent := GetRejectPercent(maxTotal, maxRejected)
	message += "Процент отмент: " + rejectedPercent

	return message
}

func CreateMessageForYesterday() string {
	message := "СТАТИСТИКА ПО ВСЕМ ЗАКАЗАМ ЗА ВЧЕРА: \n"
	message += CreateMessageForOrders([3]string{
			TARGET_ORDERS_TOTAL, 
			TARGET_ORDERS_FINISHED, 
			TARGET_ORDERS_REJECTED})

	message += "\n \n"
	message += "СТАТИСТИКА ПО КОМФОРТНЫМ ЗАКАЗАМ ЗА ВЧЕРА: \n"
	message += CreateMessageForOrders([3]string{
			TARGET_ORDERS_COMFORT_TOTAL, 
			TARGET_ORDERS_COMFORT_FINISHED, 
			TARGET_ORDERS_COMFORT_REJECTED})
	return message
}

func CreateMessageForCabs(config Config) string {
	message := "На текущий момент"
	message += "\n* свободных машин: " + strconv.Itoa(GetFreeCabsNamba(config)) + " 🆓"
	message += "\n* всего машин: " + strconv.Itoa(GetAllCabsNamba(config)) + " 🚕"
	return message
}

func SendFullInfo(config Config) {
	path    := "/tmp/drivers.png"
	orders  := CreateMessageForYesterday()
	cabs    := CreateMessageForCabs(config)
	date    := GetDayBeforeInFormat(time.Now())
	ConnectTelegramAndSendMessage(orders, config)
	checkErr(GetPicAboutCabs(date, path, config))
	ConnectTelegramAndSendPic(path, "Распределение машин за вчера", config)
	ConnectTelegramAndSendMessage(cabs, config)
}

func SendCabsInfo(config Config) {
	cabs := CreateMessageForCabs(config)
	ConnectTelegramAndSendMessage(cabs, config)
}

func main() {
	config.Fill("./config", "yml")
	log.Printf("scheduler for TFYD: %s", config.TimeForYesterdayData)
	log.Printf("scheduler for TFEC: %s", config.TimeForEveningCabs)
	jobMorning := func() {
		log.Print("morning\n")
		SendFullInfo(config)
	}
	jobEvening := func() {
		log.Print("evening\n")
		SendCabsInfo(config)
	}
	_, err := scheduler.Every().Day().At(config.TimeForYesterdayData).Run(jobMorning)
	checkErr(err)
	_, err = scheduler.Every().Day().At(config.TimeForEveningCabs).Run(jobEvening)
	checkErr(err)
	runtime.Goexit()
}

func exe(cmdName string, cmdArgs []string) string {
	var cmdOut []byte
	var	err    error
	if cmdOut, err = exec.Command(cmdName, cmdArgs...).Output(); err != nil {
		fmt.Printf("git %v error %v", cmdArgs, err)
		os.Exit(1)
	}
	return string(cmdOut)
}

func checkErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}
