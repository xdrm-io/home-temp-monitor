package main

import (
	"flag"
	"fmt"
	"log"
)

type Config struct {
	BrokerAddr     string
	BrokerUser     string
	BrokerPass     string
	SubscribeTopic string
}

func ReadConfig() (*Config, error) {
	flag.Parse()
	args := flag.Args()
	if len(args) != 4 {
		return nil, fmt.Errorf("expected 4 arguments ; got %d", len(args))
	}
	config := &Config{
		BrokerAddr:     args[0],
		BrokerUser:     args[1],
		BrokerPass:     args[2],
		SubscribeTopic: args[3],
	}
	if config.BrokerAddr == "" {
		return nil, fmt.Errorf("missing mandatory environment variable %q", "BROKER_ADDR")
	}
	if config.BrokerUser == "" {
		return nil, fmt.Errorf("missing mandatory environment variable %q", "BROKER_USER")
	}
	if config.BrokerPass == "" {
		return nil, fmt.Errorf("missing mandatory environment variable %q", "BROKER_PASS")
	}
	if config.SubscribeTopic == "" {
		return nil, fmt.Errorf("missing mandatory environment variable %q", "MQTT_TOPIC")
	}
	log.Printf("[config] broker_addr: %q", config.BrokerAddr)
	log.Printf("[config] broker_user: %q", config.BrokerUser)
	log.Printf("[config] broker_pass: %q", config.BrokerPass)
	log.Printf("[config] subscribe_topic: %q", config.SubscribeTopic)
	return config, nil
}
