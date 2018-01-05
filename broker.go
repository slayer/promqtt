package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"regexp"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	inactivityLimit = 60 * time.Second
)

var (
	messagesPublished = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "promqtt_mqtt_messages_published_total",
			Help: "Number of published messages.",
		}, []string{"name", "url"})

	messagesReceived = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "promqtt_mqtt_messages_received_total",
			Help: "Number of received messages.",
		}, []string{"name", "url"})

	errors = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "promqtt_mqtt_errors_total",
			Help: "Number of errors occurred during test execution.",
		}, []string{"name", "url"})
)

// Start reconnects
func (broker *Broker) Start() {
	broker.Init()
	for {
		broker.Connect()
		for broker.client.IsConnected() {
			time.Sleep(20 * time.Second)

			timeout := time.Now().Sub(broker.lastReceivedMessageAt)
			if timeout > inactivityLimit {
				Log.Info("Inactivity timeout, restating...", "timeout", timeout)
				broker.client.Disconnect(5)
				break
			}
		}
	}
}

// Connect connects to MQTT server and subscribe to topics
func (broker *Broker) Connect() {
	qos := byte(0)

	// Initialize optional metrics with initial values to have them present from the beginning
	messagesPublished.WithLabelValues(broker.Name, broker.URL).Add(0)
	messagesReceived.WithLabelValues(broker.Name, broker.URL).Add(0)
	errors.WithLabelValues(broker.Name, broker.URL).Add(0)

	tlsconfig := NewTLSConfig(broker)

	cliendID := "promqtt"
	clientOptions := mqtt.NewClientOptions().
		SetClientID(cliendID).
		SetUsername(broker.Username).
		SetPassword(broker.Password).
		SetTLSConfig(tlsconfig).
		SetAutoReconnect(true).
		AddBroker(broker.URL)

	// Reset acivity watchdog
	broker.lastReceivedMessageAt = time.Now()
	clientOptions.SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {
		broker.lastReceivedMessageAt = time.Now()
		messagesReceived.WithLabelValues(broker.Name, broker.URL).Inc()
		queue <- rawMessage{
			broker:  broker,
			topic:   msg.Topic(),
			payload: string(msg.Payload()),
		}
	})

	broker.client = mqtt.NewClient(clientOptions)

	if token := broker.client.Connect(); token.Wait() && token.Error() != nil {
		Log.Error("Client.Connect", token.Error())
		return
	}
	Log.Info("Connected", "name", broker.Name, "url", broker.URL)

	if token := broker.client.Subscribe(broker.Topic, qos, nil); token.Wait() && token.Error() != nil {
		Log.Error("Client Subscribe", token.Error())
		return
	}
}

// Init initialized regexps
func (broker *Broker) Init() {
	if broker.Labels != "" {
		broker.labelsRegexp = regexp.MustCompile(broker.Labels)
		labelNames := broker.labelsRegexp.SubexpNames()
		if len(labelNames) > 0 {
			broker.labelNames = labelNames[1:] // [0] always empty
		}
	}
	if broker.MetricName != "" {
		broker.metricNameRegexp = regexp.MustCompile(broker.MetricName)
	}
	return
}

// NewTLSConfig returns configured `tls.Config`
// Stolen from https://github.com/shoenig/go-mqtt/blob/master/samples/ssl.go
func NewTLSConfig(broker *Broker) *tls.Config {
	// Import trusted certificates from CAChain - purely for verification - not sent to TLS server
	certpool := x509.NewCertPool()
	pemCerts, err := ioutil.ReadFile(broker.CAChain)
	if err == nil {
		certpool.AppendCertsFromPEM(pemCerts)
	}

	// Import client certificate/key pair
	// If you want the chain certs to be sent to the server, concatenate the leaf,
	//  intermediate and root into the ClientCert file
	cert, err := tls.LoadX509KeyPair(broker.ClientCert, broker.ClientKey)
	if err != nil {
		return &tls.Config{}
	}

	// Create tls.Config with desired tls properties
	return &tls.Config{
		// RootCAs = certs used to verify server cert.
		RootCAs: certpool,
		// ClientAuth = whether to request cert from server.
		// Since the server is set up for SSL, this happens
		// anyways.
		ClientAuth: tls.NoClientCert,
		// InsecureSkipVerify = verify that cert contents
		// match server. IP matches what is in cert etc.
		InsecureSkipVerify: false,
		// Certificates = list of certs client sends to server.
		Certificates: []tls.Certificate{cert},
	}
}

func (broker *Broker) reportError(message string, err error) {
	Log.Error(message, "err", err)
}

func qqq() {
	var myExp = regexp.MustCompile(`(?P<first>\d+)\.(\d+).(?P<second>\d+)`)
	match := myExp.FindStringSubmatch("1234.5678.9")
	result := make(map[string]string)
	for i, name := range myExp.SubexpNames() {
		if i != 0 {
			result[name] = match[i]
		}
	}
	fmt.Printf("by name: %s %s\n", result["first"], result["second"])
}
