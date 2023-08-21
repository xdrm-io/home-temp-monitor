package main

import (
	"context"
	"encoding/json"
	"log"
	"regexp"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// Collector collects data from the MQTT broker.
type Collector struct {
	cnf     Config
	cli     mqtt.Client
	storage Storage
}

func NewCollector(cnf Config, storage Storage) (*Collector, error) {
	return &Collector{
		cnf:     cnf,
		storage: storage,
	}, nil
}
func (c *Collector) Close() {
	c.cli.Disconnect(1000)
}

// Subscribe subscribes to the MQTT broker.
func (c *Collector) Subscribe() error {
	if c.cli == nil {
		opts := mqtt.NewClientOptions()
		opts.AddBroker(c.cnf.BrokerAddr + ":1883")
		opts.SetUsername(c.cnf.BrokerUser)
		opts.SetPassword(c.cnf.BrokerPass)
		opts.SetClientID("station")
		c.cli = mqtt.NewClient(opts)
	}
	if !c.cli.IsConnected() {
		if token := c.cli.Connect(); token.Wait() && token.Error() != nil {
			return token.Error()
		}
	}

	if token := c.cli.Subscribe(c.cnf.SubscribeTopic, 0, c.onReceive); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	return nil
}

var roomIDReg = regexp.MustCompile(`^/room/([^/]+)/env$`)

func (c *Collector) onReceive(client mqtt.Client, msg mqtt.Message) {
	log.Printf("received: [%s] %s", msg.Topic(), msg.Payload())

	matches := roomIDReg.FindStringSubmatch(msg.Topic())
	if len(matches) != 2 {
		log.Printf("error: invalid topic, cannot get room id")
		return
	}
	roomID := matches[1]

	var m Measure
	m.Room = roomID
	if err := json.Unmarshal(msg.Payload(), &m); err != nil {
		log.Printf("error: cannot read json: %v", err)
		return
	}

	if err := c.storage.Append(context.Background(), m); err != nil {
		log.Printf("error: cannot store: %v", err)
		return
	}
	log.Printf("stored")
}
