package wolverine

import (
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"net/http"
	"time"
)

var (
	urlUp = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "sample_external_url_up",
		Help: "Is the URL up (1) or down (0).",
	})
)

func init() {
	prometheus.MustRegister(urlUp)
}

func MonitorURL(url string, httpClient *http.Client) {
	const sleepSeconds = 5

	for {
		start := time.Now()
		resp, err := httpClient.Head(url)
		if err != nil {
			log.Printf("Error: %s", err)
		} else {
			duration := time.Since(start)
			log.Printf("URL: %s, Status: %d, Duration: %dms", url, resp.StatusCode, duration.Milliseconds())
		}
		time.Sleep(sleepSeconds * time.Second)
	}
}
