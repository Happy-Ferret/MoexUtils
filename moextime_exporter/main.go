package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	listenPort = "9262"
)

func init() {
	flag.StringVar(&listenPort, "L", listenPort, "Lister port")
	flag.Parse()
}

func setTradeTime() float64 {
	t := time.Now()
	if t.Weekday() == 6 || t.Weekday() == 7 {
		return 0
	}

	tradeStart := time.Date(t.Year(), t.Month(), t.Day(), 10, 0, 0, 0, t.Location())
	tradesStop := time.Date(t.Year(), t.Month(), t.Day(), 18, 45, 0, 0, t.Location())
	cliringStart := time.Date(t.Year(), t.Month(), t.Day(), 14, 0, 0, 0, t.Location())
	cliringStop := time.Date(t.Year(), t.Month(), t.Day(), 14, 5, 0, 0, t.Location())

	if t.After(cliringStart) && t.Before(cliringStop) {
		return 0
	} else if t.After(tradesStop) {
		return 0
	} else if t.Before(tradeStart) {
		return 0
	}
	return 1
}

func main() {
	tradeTime := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "times",
		Help: "times description",
	})
	prometheus.MustRegister(tradeTime)
	tradeTime.Set(setTradeTime())
	go func() {
		for {
			tradeTime.Set(setTradeTime())
			time.Sleep(time.Duration(60 * time.Second))
			log.Println(setTradeTime())
		}
	}()

	http.Handle("/metrics", prometheus.Handler())

	http.ListenAndServe(":"+listenPort, nil)

}
