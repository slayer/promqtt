package main

import (
	"flag"
	"io/ioutil"
	"os"

	log "github.com/inconshreveable/log15"

	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"gopkg.in/yaml.v2"
)

var build string

const debug = false

var (
	// Log is a main logger
	Log   log.Logger = log.Root()
	queue            = make(chan rawMessage, 100)

	configFile    = flag.String("config.file", "config.yaml", "promqtt configuration file.")
	listenAddress = flag.String("web.listen-address", ":9214", "The address to listen on for HTTP requests.")
	logLevel      = flag.String("log.level", "info", "Log level.")
)

func init() {
	prometheus.MustRegister(messagesPublished)
	prometheus.MustRegister(messagesReceived)
}

func main() {
	flag.Parse()
	yamlFile, err := ioutil.ReadFile(*configFile)
	lvl, _ := log.LvlFromString(*logLevel)
	Log.SetHandler(
		log.LvlFilterHandler(lvl,
			log.StreamHandler(os.Stderr, log.TerminalFormat())))
	if err != nil {
		Log.Crit("Error reading config file", "err", err)
	}

	config := config{}

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		Log.Crit("Error parsing config file", "err", err)
	}
	initMetrics(&config)

	Log.Info("Starting promqtt", "build", build)

	for _, broker := range config.Brokers {
		go func(broker *Broker) {
			broker.Start()
		}(broker)
	}
	go startPrometheusExporter()

	http.Handle("/metrics", prometheus.Handler())
	http.ListenAndServe(*listenAddress, nil)
}
