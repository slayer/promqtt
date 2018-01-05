package main

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

type rawMessage struct {
	broker  *Broker
	topic   string
	payload string
}

type message struct {
	rawMessage
	labels prometheus.Labels
}

var (
	counterRegexps []*regexp.Regexp = make([]*regexp.Regexp, 0)
	gaugeRegexps   []*regexp.Regexp = make([]*regexp.Regexp, 0)
)

func initMetrics(conf *config) {
	for _, rule := range conf.Rules.Counters {
		re := regexp.MustCompile(rule)
		counterRegexps = append(counterRegexps, re)
	}
	for _, rule := range conf.Rules.Gauges {
		re := regexp.MustCompile(rule)
		gaugeRegexps = append(gaugeRegexps, re)
	}
}

// newMessage builds message from rawMessage
func newMessage(raw rawMessage) *message {
	msg := message{
		rawMessage: raw,
	}
	msg.parseLabels()

	return &msg
}

func (msg *message) metricName() string {
	name := msg.topic
	if re := msg.broker.metricNameRegexp; re != nil {
		matches := re.FindStringSubmatch(name)
		if matches != nil && len(matches) > 0 {
			name = msg.broker.MetricPrefix + matches[1]
		}
	} else {
		name = msg.broker.MetricPrefix + strings.Replace(name, "/", "_", -1)
	}
	// logger.Printf("metricName: %s -> %s \n", msg.topic, name)
	return name
}

func (msg *message) parseLabels() {
	re := msg.broker.labelsRegexp
	if re != nil {
		match := re.FindStringSubmatch(msg.topic)
		if match == nil {
			return
		}
		msg.labels = make(map[string]string)
		for i, name := range re.SubexpNames() {
			if i != 0 {
				msg.labels[name] = match[i]
			}
		}
	}

}

func (msg *message) getValue() (val float64, err error) {
	val, err = strconv.ParseFloat(msg.payload, 64)
	if err != nil {
		Log.Warn("Invalid float", "payload", msg.payload, "err", err)
	}
	return
}

func (msg *message) validLabels() (ok bool) {
	if msg.broker.labelNames == nil {
		// Broker has no labels, no need to check it in message
		return true
	}
	ok = true

	// no labels in msg found, but broker have
	ok = !(msg.labels == nil || len(msg.labels) == 0)
	// msg labels has exact len as broker's labels len
	ok = ok && len(msg.labels) == len(msg.broker.labelNames)

	if !ok {
		Log.Debug("Fail to parse labels", "topic", msg.topic)
	}

	return ok
}

func (msg *message) isCounter() bool {
	for _, re := range counterRegexps {
		if re.MatchString(msg.topic) {
			return true
		}
	}
	return false
}

func (msg *message) isGauge() bool {
	for _, re := range gaugeRegexps {
		if re.MatchString(msg.topic) {
			return true
		}
	}
	return false
}

func isLooksLikeAMetric(val string) bool {
	if len(val) == 0 {
		return false
	}

	c := val[0]
	if (c >= '0' && c <= '9') || c == '.' || c == '-' || c == '+' {
		return floatRegexp.MatchString(val)
	}
	return false
}
