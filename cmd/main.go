package main

import (
	"github.com/adrian/wolverine/internal"
	"github.com/kkyr/fig"
	"log"
)

type Config struct {
	URLs []string `fig:"urls" default:"[]"`
}

func main() {
	// load URLs to monitor from config file
	var cfg Config
	err := fig.Load(&cfg, fig.Dirs("config"))
	if err != nil {
		log.Fatal(err)
	}

	wolverine.Monitor(cfg.URLs)
}
