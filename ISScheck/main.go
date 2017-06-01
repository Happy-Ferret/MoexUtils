package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"time"

	moexlib "github.com/agareev/MoexLib/monitoring"
	config "github.com/agareev/MoexLib/other"
	"github.com/jasonlvhit/gocron"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
)

/*
{
"marketdata": {
        "metadata": {
                "UPDATETIME": {"type": "time", "bytes": 10, "max_size": 0}
        },
        "columns": ["UPDATETIME"],
        "data": [
                ["18:01:08"]
        ]
}}
*/

// Request JSON description
type Request struct {
	Marketdata struct {
		Columns []string   `json:"columns"`
		Data    [][]string `json:"data"`
	} `json:"marketdata"`
	Trade struct {
		Columns []string   `json:"columns"`
		Data    [][]string `json:"data"`
	} `json:"trades"`
}

var (
	configuration   config.Config
	issURL          = "http://iss.moex.com/iss"
	isMOCK          bool
	gatewayURL      = "http://gturl:9091"
	throughputGuage = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "iss",
		Help: "help iss",
	})
)

func init() {
	configuration = config.ReadConfig("config.json")
	// add flag parse, add mock function
}

// GetDelta return string delta time
func getDelta(lastDealTime string) float64 {
	LastDeal := moexlib.StringTime2UnixTime(lastDealTime)
	delta := math.Abs(time.Now().Sub(LastDeal).Seconds())
	return delta
}

func randNum() string {
	return fmt.Sprintf("&rand=%v", rand.Intn(1000))
}

func urlReturn(engine, market, typeOfCheck string) string {
	// engine - stock, futures, currency, stock
	// market - index, forts, selt, shares
	var parturl string
	if typeOfCheck == "marketdata" {
		parturl = "/securities.json?iss.only=marketdata&sort_order=desc&first=1"
		if market == "shares" {
			// FIXME временный work around пока не исправим запрос в ИСС
			parturl += "&marketdata.columns=TIME&sort_column=TIME"
		} else {
			parturl += "&marketdata.columns=UPDATETIME&sort_column=UPDATETIME"
		}
		if market == "index" {
			// Мониториим только realtime индексы.
			// Выбраны основные индексы которыи приходят из разных считалок
			parturl += "&securities=MICEXINDEXCF,RTSI,MICEXBMI,RTSSTD,RVI"
		}
	} else if typeOfCheck == "trades" {
		parturl = "/trades.json?reversed=1&limit=1&iss.only=trades&trades.columns=TRADETIME"
	} else {
		log.Fatal("unknown type of check")
		return "unknown type of check"
	}
	return issURL + "/engines/" + engine + "/markets/" + market + parturl + randNum()
}

// TODO split on 2 functions
func getURL(url string) string {
	var input Request
	var output string

	// log.Println(url)
	json.Unmarshal(moexlib.GetAllContents(url), &input)
	if input.Marketdata.Columns == nil {
		json.Unmarshal(moexlib.GetAllContents(url), &input)
		output = input.Trade.Data[0][0]
		return output
	}
	output = input.Marketdata.Data[0][0]
	return output
}

func execute() {
	engines := [][2]string{
		{"stock", "shares"},
		{"currency", "selt"},
		{"futures", "forts"},
		{"stock", "index"},
	}

	checks := []string{"marketdata", "trades"}
	for _, typeOfCheck := range checks {
		for _, marketInfo := range engines {
			engine, market := marketInfo[0], marketInfo[1]
			url := urlReturn(engine, market, typeOfCheck)
			diff := getDelta(getURL(url))
			delta := fmt.Sprintf("%v", diff)
			// fmt.Println(engine+"--"+market, delta, url)
			log.Println("Got " + engine + "_" + market + " delta: " + delta)
			throughputGuage = prometheus.NewGauge(prometheus.GaugeOpts{
				Name: "iss_" + typeOfCheck + "_" + engine + "_" + market,
				Help: "iss." + typeOfCheck + "." + engine + "." + market + " in seconds",
			})
			throughputGuage.Set(diff)
			if err := push.Collectors("throughput_job", push.HostnameGroupingKey(), gatewayURL, throughputGuage); err != nil {
				fmt.Println("Could not push completion time to Pushgateway:", err)
			}
		}
	}
	log.Println("Checked all data")
}

func main() {
	s := gocron.NewScheduler()
	s.Every(5).Seconds().Do(execute)
	<-s.Start()
}
