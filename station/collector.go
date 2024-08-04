package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"regexp"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/xdrm-io/home-temp-monitor/storage"
)

// Collector collects data from the MQTT broker.
type Collector struct {
	cnf     Config
	cli     mqtt.Client
	storage storage.Storage
}

func NewCollector(cnf Config, s storage.Storage) (*Collector, error) {
	return &Collector{
		cnf:     cnf,
		storage: s,
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

	var m storage.Measure
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

	c.notifyHass(context.Background(), roomID, m)
}

// notifyHass notifies Home Assistant of the new measure.
func (c *Collector) notifyHass(ctx context.Context, room string, m storage.Measure) {
	// only take the current measure into account, not cached ones
	if m.OffsetSec != 0 {
		return
	}

	type Config struct {
		DeviceClass       string `json:"device_class"`
		UniqueID          string `json:"unique_id"`
		Name              string `json:"name"`
		StateTopic        string `json:"state_topic"`
		UnitOfMeasurement string `json:"unit_of_measurement"`
		ValueTemplate     string `json:"value_template"`
	}

	var (
		topicT = "homeassistant/sensor/" + room + "_t"
		topicH = "homeassistant/sensor/" + room + "_h"

		t = float64(m.Temperature) / 10.
		h = float64(m.Humidity) / 10.
	)

	tconfig, err := json.Marshal(Config{
		DeviceClass:       "temperature",
		Name:              room + " temperature",
		UniqueID:          "sensor." + room + "_t",
		StateTopic:        topicT + "/state",
		UnitOfMeasurement: "°C",
		ValueTemplate:     "{{ value }}",
	})
	if err != nil {
		log.Printf("hass notify: encode temperature config: %v", err)
		return
	}

	hconfig, err := json.Marshal(Config{
		DeviceClass:       "humidity",
		Name:              room + " humidity",
		UniqueID:          "sensor." + room + "_h",
		StateTopic:        topicH + "/state",
		UnitOfMeasurement: "%",
		ValueTemplate:     "{{ value }}",
	})
	if err != nil {
		log.Printf("hass notify: encode humidity config: %v", err)
		return
	}

	// config
	if tok := c.cli.Publish(topicT+"/config", 0, false, tconfig); tok.Error() != nil {
		log.Printf("hass notify: publish tconfig: %v", tok.Error())
	}
	if tok := c.cli.Publish(topicH+"/config", 0, false, hconfig); tok.Error() != nil {
		log.Printf("hass notify: publish hconfig: %v", tok.Error())
	}

	// state
	tstate := fmt.Sprintf(`%.1f`, t)
	if tok := c.cli.Publish(topicT+"/state", 0, false, tstate); tok.Error() != nil {
		log.Printf("hass notify: publish tstate: %v", tok.Error())
	}
	hstate := fmt.Sprintf(`%.1f`, h)
	if tok := c.cli.Publish(topicH+"/state", 0, false, hstate); tok.Error() != nil {
		log.Printf("hass notify: publish hstate: %v", tok.Error())
	}
}
