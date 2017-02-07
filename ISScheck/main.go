package main

import (
	"encoding/json"
	"fmt"

	moexlib "github.com/agareev/moexLib/monitoring"
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
		"stock":    "shares",
		"currency": "selt",
		"futures":  "forts",
		// "stock":    "index",
	}

	for x, y := range urls {
		url := "http://moex.com/iss/engines/" + x + "/markets/" + y + "/securities.json?iss.only=marketdata&sort_column=updatetime&sort_order=desc&first=1&marketdata.columns=UPDATETIME"

		var output Marketdata

		contents := moexlib.GetAllContents(url)
		json.Unmarshal(contents, &output)
		// fmt.Println(output.MarketData.Columns[0])
		// fmt.Println(output.MarketData.Data[0][0])
		diff := moexlib.GetDelta(output.MarketData.Data[0][0])
		delta := fmt.Sprintf("%v", diff)
		fmt.Println(delta, output.MarketData.Data[0][0])
		ok := moexlib.Send2Graphite(delta, "iss.marketdata."+x+"."+y, "127.0.0.1", 32768)
		if ok == false {
			fmt.Println("good")
		}
	}
	for x, y := range urls {
		url := "http://moex.com/iss/engines/" + x + "/markets/" + y + "/trades.json?reversed=1&limit=1&iss.only=trades&trades.columns=TRADETIME"

		var output TradesPage

		contents := moexlib.GetAllContents(url)
		json.Unmarshal(contents, &output)
		// fmt.Println(output.MarketData.Columns[0])
		// fmt.Println(output.MarketData.Data[0][0])
		diff := moexlib.GetDelta(output.Trades.Data[0][0])
		delta := fmt.Sprintf("%v", diff)
		fmt.Println(delta, output.Trades.Data[0][0])
		ok := moexlib.Send2Graphite(delta, "iss.trades."+x+"."+y, "127.0.0.1", 32768)
		if ok == false {
			fmt.Println("good")
		}
	}
}
