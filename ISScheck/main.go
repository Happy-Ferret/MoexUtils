package main

import (
	"encoding/json"
	"fmt"
	"log"

	moexlib "github.com/agareev/MoexLib/monitoring"
	config "github.com/agareev/MoexLib/other"
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

// Marketdata description
type Marketdata struct {
	MarketData struct {
		Columns []string   `json:"columns"`
		Data    [][]string `json:"data"`
	} `json:"marketdata"`
}

// TradesPage description
type TradesPage struct {
	Trades struct {
		Columns []string   `json:"columns"`
		Data    [][]string `json:"data"`
	} `json:"trades"`
}

var url string

func main() {
	urls := map[string]string{
		"shares": "stock",
		"selt":   "currency",
		"forts":  "futures",
		"index":  "stock",
	}
	var configuration config.Config
	configuration = config.ReadConfig("config.json")

	for x, y := range urls {
		url := "http://moex.com/iss/engines/" + y + "/markets/" + x + "/securities.json?iss.only=marketdata&sort_column=updatetime&sort_order=desc&first=1&marketdata.columns=UPDATETIME"

		var output Marketdata

		contents := moexlib.GetAllContents(url)
		json.Unmarshal(contents, &output)
		// fmt.Println(output.MarketData.Columns[0])
		// fmt.Println(output.MarketData.Data[0][0])
		diff := moexlib.GetDelta(output.MarketData.Data[0][0])
		delta := fmt.Sprintf("%v", diff)
		fmt.Println(delta, output.MarketData.Data[0][0])
		ok := moexlib.Send2Graphite(delta, "iss.marketdata."+y+"."+x, configuration.Server.IP, configuration.Server.Port)
		if ok == false {
			fmt.Println("good")
		}
	}
	for x, y := range urls {
		url := "http://moex.com/iss/engines/" + y + "/markets/" + x + "/trades.json?reversed=1&limit=1&iss.only=trades&trades.columns=TRADETIME"

		var output TradesPage

		contents := moexlib.GetAllContents(url)
		json.Unmarshal(contents, &output)
		// fmt.Println(output.MarketData.Columns[0])
		// fmt.Println(output.MarketData.Data[0][0])
		diff := moexlib.GetDelta(output.Trades.Data[0][0])
		delta := fmt.Sprintf("%v", diff)
		// fmt.Println(delta, output.Trades.Data[0][0])
		err := moexlib.Send2Graphite(delta, "iss.trades."+y+"."+x, "172.22.192.91", 2003)
		if err != false {
			log.Fatal(err)
		}
	}
}
