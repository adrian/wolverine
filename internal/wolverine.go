package wolverine

import (
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"net/http"
	"time"
)

var (
	urlUpMetric = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "sample_external_url_up",
		Help: "Is the URL up (1) or down (0).",
	},
		[]string{"url"})

	responseTimeMetric = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "sample_external_url_response_ms",
		Help:    "Response time in milliseconds.",
		Buckets: []float64{50, 100, 250, 500, 750, 1000},
	},
		[]string{"url"},
	)
)

func init() {
	prometheus.MustRegister(urlUpMetric)
	prometheus.MustRegister(responseTimeMetric)
}

func MonitorURL(url string, httpClient *http.Client) (error, *http.Response) {
	start := time.Now()
	resp, err := httpClient.Head(url)
	if err != nil {
		log.Printf("Error: %v", err)
		urlUpMetric.With(prometheus.Labels{"url": url}).Set(1)
	} else {
		duration := time.Since(start)
		up := 0
		if resp.StatusCode == 200 {
			up = 1
		}
		urlUpMetric.With(prometheus.Labels{"url": url}).Set(float64(up))
		responseTimeMetric.With(prometheus.Labels{"url": url}).Observe(float64(duration.Milliseconds()))
		log.Printf("URL: %s, Status: %d, Duration: %dms", url, resp.StatusCode, duration.Milliseconds())
	}
	return err, resp
}
