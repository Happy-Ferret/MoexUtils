package mvp

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	log.Println("I am started")
}

func main() {
	metrics := make(map[string]prometheus.Gauge)
	checks := []string{"FirstMetric", "SecondMetric"}
	// registry := prometheus.NewRegistry()

	for _, name := range checks {
		metrics[name] = prometheus.NewGauge(prometheus.GaugeOpts{
			Name:        name,
			Help:        name + " in seconds",
			ConstLabels: prometheus.Labels{"stream": "opisanie"},
		})
		metrics[name].Set(11)
	}
	for key := range metrics {
		prometheus.MustRegister(metrics[key])
	}
	// metric.Set(11)
	// prometheus.MustRegister(metric)

	http.Handle("/metrics", prometheus.Handler())
	http.ListenAndServe(":8080", nil)
}
