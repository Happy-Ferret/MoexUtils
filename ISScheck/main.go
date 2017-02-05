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

func (c *Marketdata) delta() string {
	return c.Space.Data[0][0]
}

func (c *TradesPage) delta() string {
	return c.Space.Data[0][0]
}

type ddddd interface {
	delta() string
}

var url string
var configuration config.Config

func check(parturl1, parturl2, checkUnit string, urls map[string]string) {
	if checkUnit == "marketdata" {
		var output Marketdata
		for x, y := range urls {
			url := "http://iss.moex.com/iss/engines/" + y + parturl1 + x + parturl2
			contents := moexlib.GetAllContents(url)
			json.Unmarshal(contents, &output)
			diff := moexlib.GetDelta(output.delta())
			delta := fmt.Sprintf("%v", diff)
			// fmt.Println(delta, output.MarketData.Data[0][0])
			err := moexlib.Send2Graphite(delta, "iss."+checkUnit+"."+y+"."+x, configuration.Server.IP, configuration.Server.Port)
			if err != false {
				log.Fatal(err)
			}
		}
	} else {
		var output TradesPage
		for x, y := range urls {
			url := "http://iss.moex.com/iss/engines/" + y + "/markets/" + x + "/trades.json?reversed=1&limit=1&iss.only=trades&trades.columns=TRADETIME"

			contents := moexlib.GetAllContents(url)
			json.Unmarshal(contents, &output)
			diff := moexlib.GetDelta(output.delta())
			delta := fmt.Sprintf("%v", diff)
			// fmt.Println(delta, output.Trades.Data[0][0])
			err := moexlib.Send2Graphite(delta, "iss.trades."+y+"."+x, "127.0.0.1", 2003)
			if err != false {
				log.Fatal(err)
			}
		}
	}

}

func getState(urls map[string]string, checkUnit string) {
	if checkUnit == "marketdata" {
		parturl1 := "/markets/"
		parturl2 := "/securities.json?iss.only=marketdata&sort_column=UPDATETIME&sort_order=desc&first=1&marketdata.columns=UPDATETIME"
		check(parturl1, parturl2, checkUnit, urls)
	}

}

func main() {
	urls := map[string]string{
		"shares": "stock",
		"selt":   "currency",
		"forts":  "futures",
		"index":  "stock",
	}

	configuration = config.ReadConfig("config.json")
	getState(urls, "trades")
	getState(urls, "marketdata")

}
