package main

import (
	"log"

	config "github.com/agareev/MoexLib/other"
)

var configuration config.Config

// Marketdata description
type Marketdata struct {
	Space struct {
		Columns []string   `json:"columns"`
		Data    [][]string `json:"data"`
	} `json:"marketdata"`
}

// TradesPage description
type TradesPage struct {
	Space struct {
		Columns []string   `json:"columns"`
		Data    [][]string `json:"data"`
	} `json:"trades"`
}

func urlReturn(engine, market, typeOfCheck string) string {
	// engine - stock, futures, currency, stock
	// market - index, forts, selt, shares
	if typeOfCheck == "marketdata" {
		parturl := "/securities.json?iss.only=marketdata&sort_column=UPDATETIME&sort_order=desc&first=1&marketdata.columns=UPDATETIME"
		url := "http://iss.moex.com/iss/engines/" + market + "/markets/" + engine + parturl
		return url
	} else if typeOfCheck == "trades" {
		parturl := "/trades.json?reversed=1&limit=1&iss.only=trades&trades.columns=TRADETIME"
		url := "http://iss.moex.com/iss/engines/" + market + "/markets/" + engine + parturl
		return url
	} else {
		log.Fatal("bad")
		return "bad"
	}
}

func getURL(url string) {
		log.Println(url)
		log.Println("----")
	}

func main() {
	urls := map[string]string{
		"shares": "stock",
		"selt":   "currency",
		"forts":  "futures",
		"index":  "stock",
	}

	configuration = config.ReadConfig("config.json")
	checks := [2]string{"marketdata", "trades"}
	for _, typeOfCheck := range checks {
	for engine, market := range urls {
		url := urlReturn(engine, market,typeOfCheck)
		getURL(url)
	}
	}

}
