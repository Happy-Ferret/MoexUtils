package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

/*
 {
    "module_id": "authen-prod-prod-infocx11",
    "class": "ICX Authenticator v.3.0.1.15",
    "creation": "2017-06-05 08:45:58",
    "activity": "2017-07-03 15:50:10",
    "color": "green",
    "roles": "[author]",
    "timeToKill": "null"
  }
*/

type request struct {
	ModuleID string `json:"module_id"`
	Color    string `json:"color"`
}

const (
	uri = "http://172.19.132.2:8090/status"
)

var (
	metrics = make(map[string]prometheus.Gauge)
)

func getContent() (input []request) {
	response, err := http.Get(uri)
	if err != nil {
		log.Fatal(err)
	} else {
		defer response.Body.Close()
	}
	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal(content, &input)
	return input
}

func register() {
	// input := getContent()
	for _, info := range getContent() {
		replacedName := strings.Replace(info.ModuleID, "-", "_", -1)
		replacedName = strings.Replace(replacedName, ".", "_", -1)
		metricName := "infocx_" + replacedName
		log.Println(metricName)
		metrics[metricName] = prometheus.NewGauge(prometheus.GaugeOpts{
			Name: metricName,
			Help: metricName + " state",
			// ConstLabels: prometheus.Labels{"stream": metricName},
		})
		// log.Println(replacedName + " will register")
		prometheus.MustRegister(metrics[metricName])
		// log.Println(metricName + " registered")
	}
	//
}

func execute() {
	//
}

func main() {
	log.Println("I'm started")
	register()

	http.Handle("/metrics", prometheus.Handler())
	http.ListenAndServe(":9261", nil)
	// for _, i := range input {
	//   log.Println(i.ModuleID, i.Color)
	// }
}
