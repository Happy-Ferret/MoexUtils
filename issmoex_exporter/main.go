package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/http"
	"time"

	moexlib "github.com/agareev/MoexLib/monitoring"
	"github.com/jasonlvhit/gocron"
	"github.com/prometheus/client_golang/prometheus"
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
	issURL            = "http://iss.moex.com/iss"
	isMOCK     bool   = false
	debug      bool   = false
	checktime  uint64 = 15
	ListenPort        = ":9260"
	metrics           = make(map[string]prometheus.Gauge)
	engines           = [][2]string{
		{"stock", "shares"},
		{"currency", "selt"},
		{"futures", "forts"},
		{"stock", "index"},
	}

	checks = []string{"marketdata", "trades"}
)

func init() {
	flag.StringVar(&ListenPort, "L", ListenPort, "Lister port, defaul :9260")
	flag.Uint64Var(&checktime, "S", checktime, "Check interval")
	flag.BoolVar(&debug, "v", debug, "enable debug")
	flag.Parse()
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
	return issURL + "/engines/" + engine + "/markets/" + market + parturl
}

// TODO split on 2 functions
func getURL(url string) string {
	var input Request
	var output string

	json.Unmarshal(moexlib.GetAllContents(url), &input)
	if input.Marketdata.Columns == nil {
		json.Unmarshal(moexlib.GetAllContents(url), &input)
		// FIXME workarround for empty array
		if len(input.Trade.Data) != 0 {
			output = input.Trade.Data[0][0]
			return output
		} else {
			return ""
		}
	}
	output = input.Marketdata.Data[0][0]
	return output
}

func registerISScheck() {
	for _, typeOfCheck := range checks {
		for _, marketInfo := range engines {
			engine, market := marketInfo[0], marketInfo[1]
			metricName := "iss_" + typeOfCheck + "_" + engine + "_" + market
			metrics[metricName] = prometheus.NewGauge(prometheus.GaugeOpts{
				Name:        metricName,
				Help:        metricName + " in seconds",
				ConstLabels: prometheus.Labels{"stream": metricName},
			})
			prometheus.MustRegister(metrics[metricName])
			log.Println(metricName + " registered")
		}
	}
}

func executeISScheck() {
	for _, typeOfCheck := range checks {
		for _, marketInfo := range engines {

			engine, market := marketInfo[0], marketInfo[1]
			url := urlReturn(engine, market, typeOfCheck) + randNum()
			if debug == true {
				log.Println(url)
			}
			diff := getDelta(getURL(url))
			log.Println("Got "+typeOfCheck+"_"+engine+"_"+market+" delta: ", diff)
			metricName := "iss_" + typeOfCheck + "_" + engine + "_" + market
			metrics[metricName].Set(diff)
		}
	}

	log.Println("Checked all data")
}

func cron() {
	registerISScheck()
	registerUSDEURcheck()

	s := gocron.NewScheduler()
	s.Every(checktime).Seconds().Do(executeISScheck)
	s.Every(checktime).Seconds().Do(executeUSDEURcheck)
	<-s.Start()

}

func main() {
	go cron()

	http.Handle("/metrics", prometheus.Handler())
	http.ListenAndServe(ListenPort, nil)
}
