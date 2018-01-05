package main

import (
	"fmt"
	"regexp"
	"unicode"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	counters    map[string]*prometheus.CounterVec = make(map[string]*prometheus.CounterVec)
	gauges      map[string]*prometheus.GaugeVec   = make(map[string]*prometheus.GaugeVec)
	floatRegexp                                   = regexp.MustCompile(`^[-+]?[0-9]*\.?[0-9]+([eE][-+]?[0-9]+)?$`)
)

func startPrometheusExporter() {
	for raw := range queue {
		if isLooksLikeAMetric(raw.payload) {
			msg := NewMessage(raw)
			updateMetric(msg)
		}
	}
}

func updateMetric(msg *message) {
	if msg.isCounter() {
		updateCounter(msg)
	} else if msg.isGauge() {
		updateGauge(msg)
	} else {
		// Log.Debug("Looks like a metric but not exists in rules",
		// 	"name", msg.metricName(), "payload", msg.payload, "topic", msg.topic)
	}
}

func updateCounter(msg *message) {
	val, err := msg.getValue()
	if err != nil {
		return
	}
	if !msg.validLabels() {
		return
	}
	counter := getCounter(msg)
	if counter == nil {
		return
	}
	if debug {
		Log.Debug("update counter", "topic", msg.topic,
			"name", msg.metricName(), "labels", msg.labels,
			"val", val)
	}
	counter.With(msg.labels).Add(val)
}

func updateGauge(msg *message) {
	val, err := msg.getValue()
	if err != nil {
		return
	}
	if !msg.validLabels() {
		return
	}
	gauge := getGauge(msg)
	if gauge == nil {
		return
	}
	if debug {
		Log.Debug("update gauge", "topic", msg.topic,
			"name", msg.metricName(), "labels", msg.labels,
			"val", val)
	}
	gauge.With(msg.labels).Set(val)
}

func getCounter(msg *message) *prometheus.CounterVec {
	name := msg.metricName()
	if counter := counters[name]; counter != nil {
		return counter
	}
	counter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: name,
		Help: fmt.Sprintf("Counter for %s metric", name),
	}, msg.broker.labelNames)

	if err := prometheus.Register(counter); err != nil {
		Log.Error("Cannot register", "name", name, "err", err)
		return nil
	}
	counters[name] = counter
	return counter
}

func getGauge(msg *message) *prometheus.GaugeVec {
	name := msg.metricName()
	if gauge := gauges[name]; gauge != nil {
		return gauge
	}
	gauge := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: name,
		Help: fmt.Sprintf("Gauge for %s metric", name),
	}, msg.broker.labelNames)

	if err := prometheus.Register(gauge); err != nil {
		Log.Error("Cannot register", "name", name, "err", err)
		return nil
	}
	gauges[name] = gauge
	return gauge
}

// IsASCIIPrintable checks if s is ascii and printable,
// aka doesn't include tab, backspace, etc.
func IsASCIIPrintable(s string) bool {
	for _, r := range s {
		if r > unicode.MaxASCII || !unicode.IsPrint(r) {
			return false
		}
	}
	return true
}
