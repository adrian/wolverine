package main

import (
	"fmt"
	"github.com/adrian/wolverine/internal"
	"github.com/kkyr/fig"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"time"
)

type Config struct {
	URLs []string `fig:"urls" default:"[]"`
}

var BuildVersion = "dev"

func main() {
	const metricsPort = 2112
	const metricsPath = "/metrics"
	const requestTimeoutSeconds = 5

	log.Printf("Starting wolverine version: %s", BuildVersion)
	log.Printf("Starting prometheus metrics endpoint on port: %d, path: %s",
		metricsPort, metricsPath)
	http.Handle(metricsPath, promhttp.Handler())
	go http.ListenAndServe(fmt.Sprintf(":%d", metricsPort), nil)

	// load URLs to monitor from config file
	var cfg Config
	err := fig.Load(&cfg, fig.Dirs("config"))
	if err != nil {
		log.Fatal(err)
	}

	// monitor each URL in a goroutine
	httpClient := &http.Client{Timeout: requestTimeoutSeconds * time.Second}
	for _, url := range cfg.URLs {
		go monitorURL(url, httpClient)
	}

	// wait indefinately
	select {}
}

func monitorURL(url string, httpClient *http.Client) {
	const sleepSeconds = 5
	for {
		wolverine.MonitorURL(url, httpClient)
		time.Sleep(sleepSeconds * time.Second)
	}
}
