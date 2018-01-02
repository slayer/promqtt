package main

import (
	"flag"
	"io/ioutil"

	log "github.com/inconshreveable/log15"

	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"gopkg.in/yaml.v2"
)

var build string

const debug = false

var (
	Log   log.Logger = log.Root()
	queue            = make(chan rawMessage, 100)

	probeDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "promqtt_mqtt_duration_seconds",
			Help: "Time taken to execute probe.",
		}, []string{"name", "url"})

	configFile    = flag.String("config.file", "config.yaml", "promqtt configuration file.")
	listenAddress = flag.String("web.listen-address", ":9214", "The address to listen on for HTTP requests.")
)

func init() {
	prometheus.MustRegister(probeDuration)
	prometheus.MustRegister(messagesPublished)
	prometheus.MustRegister(messagesReceived)
}

func main() {
	flag.Parse()
	yamlFile, err := ioutil.ReadFile(*configFile)

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
