package main

import (
	"encoding/json"
	"log"
	"strconv"

	m "ops.MonitoringScripts/monitoringlibs"
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

func getCurrency(content []byte) string {
	var output ResponseISS
	var lastI int
	json.Unmarshal(content, &output)
	lastI = getIndex(output.Marketdata.Columns)
	return strconv.FormatFloat(output.Marketdata.Data[0][lastI], 'f', -1, 64)
}

func registerUSDEURcheck() {
	log.Println("registered")
}

func executeUSDEURcheck() {
	secs := make(map[string]string)
	secs["usd"] = usd
	secs["eur"] = eur

	for key, sec := range secs {
		content := m.GetAllContents(issURL + sec)
		delta := getCurrency(content)
		log.Println("executed", key, delta)
	}
}
