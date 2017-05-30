package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"regexp"
	"strings"
	// . "github.com/blacked/go-zabbix"
	"github.com/jasonlvhit/gocron"
	"github.com/prometheus/client_golang/prometheus"
	m "ops.MonitoringScripts/monitoringlibs"
)

// plazaResponse - json
type plazaResponse struct {
	PlazaThreads []struct {
		Connection string `json:"connection"`
		Status     string `json:"status"`
		Listeners  []struct {
			Stream string `json:"stream"`
			Mode   string `json:"mode"`
			Status string `json:"status"`
			Tables []struct {
				ID    string      `json:"id"`
				Delay int         `json:"delay"`
				Queue int         `json:"queue"`
				Ts    interface{} `json:"ts"`
			} `json:"tables"`
		} `json:"listeners"`
	} `json:"plaza_threads"`
}

var (
	metrics            = make(map[string]prometheus.Gauge)
	astsAdress         = "127.0.0.1"
	plazaAdress        = "127.0.0.1"
	astsPort           = "8888"
	plazaPort          = "9999"
	listenPort         = ":9261"
	checktime   uint64 = 15
)

func init() {
	flag.StringVar(&listenPort, "L", listenPort, "Lister port")
	flag.StringVar(&astsAdress, "as", astsAdress, "Asts-Gate ip address")
	flag.StringVar(&astsPort, "ap", astsPort, "Asts-Gate port")
	flag.StringVar(&plazaAdress, "pl", plazaAdress, "Plaza-Gate ip address")
	flag.StringVar(&plazaPort, "pp", plazaPort, "Plaza-Gate port")
	flag.Uint64Var(&checktime, "c", checktime, "frequency checks")
	flag.Parse()
}

func urlReturn(typeOfcheck string) string {
	if typeOfcheck == "asts" {
		return "http://" + astsAdress + ":" + astsPort
	}
	return "http://" + plazaAdress + ":" + plazaPort + "/status.json"
}

// AstsCheck - check asts gate functions
func astsCheck() {
	ServerURL := urlReturn("asts")
	contents := strings.Split(string(m.GetAllContents(ServerURL)), "\n")
	regex := ">(?P<first>stock|currency|interventions)<.*>(?P<second>alive|dead)</td>"
	streamRegex := regexp.MustCompile(regex)
	// TODO Fixme please
	for _, name := range contents {
		match := streamRegex.FindStringSubmatch(name)
		result := make(map[string]string)
		if match == nil {
			continue
		}

		for i, name := range streamRegex.SubexpNames() {
			if (i != 0) && (match[i] != "") {
				result[name] = match[i]
				if result["second"] == "" {
					continue
				}
				state := 0
				if result["second"] != "alive" {
					state = 0
				} else {
					state = 1
				}
				metric := prometheus.NewGauge(
					prometheus.GaugeOpts{
						Name: result["first"],
						Help: result["first"] + " state",
						// ConstLabels: prometheus.Labels{"stream": metricName},
					})
				prometheus.Register(metric)
				if state == 1 {
					metric.Inc()
				} else {
					metric.Set(0)
				}

				log.Println(result["first"],
					state)
			}
		}
	}
}

// PlazaCheck - check plaza gate functions
func plazaCheck() {
	var output plazaResponse
	ServerURL := urlReturn("plaza")
	contents := m.GetAllContents(ServerURL)

	json.Unmarshal(contents, &output)
	for _, content := range output.PlazaThreads[0].Listeners {
		var streamRegex = regexp.MustCompile(`//(.+);`)
		streamString := streamRegex.FindStringSubmatch(content.Stream)[1]
		state := 0

		if content.Status != "Online" {
			state = 0
		} else {
			state = 1
		}
		metric := prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: streamString,
				Help: streamString + " state",
				// ConstLabels: prometheus.Labels{"stream": metricName},
			})
		prometheus.Register(metric)
		if state == 1 {
			metric.Inc()
		} else {
			metric.Set(0)
		}

		log.Println(streamString,
			state)
	}
}

func cron() {
	s := gocron.NewScheduler()
	s.Every(checktime).Seconds().Do(astsCheck)
	s.Every(checktime).Seconds().Do(plazaCheck)
	<-s.Start()
}

func main() {
	go cron()

	http.Handle("/metrics", prometheus.Handler())
	http.ListenAndServe(listenPort, nil)
	// TODO add restart add starttime, queue
}
