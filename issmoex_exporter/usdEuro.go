package main

import (
	"encoding/json"
	"log"

	moexlib "github.com/agareev/MoexLib/monitoring"
	"github.com/prometheus/client_golang/prometheus"
)

const url = "/engines/currency/markets/selt/boards/CETS/securities/"
const usd = "USD000UTSTOM.json?iss.meta=off&iss.only=marketdata"
const eur = "EUR_RUB__TOM.json?iss.meta=off&iss.only=marketdata"

// ResponseISS is a json ResponseISS from iss.moex.com
type ResponseISS struct {
	Marketdata struct {
		Columns []string    `json:"columns"`
		Data    [][]float64 `json:"data"`
	} `json:"marketdata"`
}

func getIndex(array []string) int {
	for i, name := range array {
		if name == "LAST" {
			return i
		}
	}
	return 0
}

func getCurrency(content []byte) float64 {
	var output ResponseISS
	var lastI int
	json.Unmarshal(content, &output)
	lastI = getIndex(output.Marketdata.Columns)
	return output.Marketdata.Data[0][lastI]
}

func registerMetric(id, name string) {
	metrics[id] = prometheus.NewGauge(prometheus.GaugeOpts{
		Name:        name,
		Help:        name + " in rubles",
		ConstLabels: prometheus.Labels{"stream": name},
	})
	prometheus.MustRegister(metrics[id])
	log.Println(id + " registered")
}

func registerUSDEURcheck() {
	registerMetric("metricUSD", "iss_usd")
	registerMetric("metricEUR", "iss_eur")
}

func executeUSDEURcheck() {
	secs := make(map[string]string)
	secs["usd"] = usd
	secs["eur"] = eur

	for key, sec := range secs {
		if debug == true {
			log.Println(issURL + url + sec)
		}
		content := moexlib.GetAllContents(issURL + url + sec)
		delta := getCurrency(content)
		if key == "usd" {
			metrics["metricUSD"].Set(delta)
		} else {
			metrics["metricEUR"].Set(delta)
		}
		log.Println("Got", key, "delta:", delta)
	}
}
