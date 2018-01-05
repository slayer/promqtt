package main

import (
	"regexp"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type config struct {
	Brokers []*Broker   `yaml:"brokers"`
	Rules   configRules `yaml:"rules"`
}

type configRules struct {
	Counters []string `yaml:"counters"`
	Gauges   []string `yaml:"gauges"`
}

// Broker is a struct for MQTT->Prom converter
type Broker struct {
	Name  string `yaml:"name"`
	URL   string `yaml:"url"`
	Topic string `yaml:"topic"`
	// ClientPrefix string `yaml:"client_prefix"`
	Username     string `yaml:"username"`
	Password     string `yaml:"password"`
	ClientCert   string `yaml:"client_cert"`
	ClientKey    string `yaml:"client_key"`
	CAChain      string `yaml:"ca_chain"`
	MetricName   string `yaml:"metric_name"`
	MetricPrefix string `yaml:"metric_prefix"`
	Labels       string `yaml:"labels"`
	// internal fields
	labelsRegexp          *regexp.Regexp
	labelNames            []string
	metricNameRegexp      *regexp.Regexp
	client                mqtt.Client
	lastReceivedMessageAt time.Time
}
