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
	DBPath         string
}

func ReadConfig() (*Config, error) {
	flag.Parse()
	args := flag.Args()
	if len(args) != 5 {
		return nil, fmt.Errorf("expected 5 arguments ; got %d", len(args))
	}
	config := &Config{
		BrokerAddr:     args[0],
		BrokerUser:     args[1],
		BrokerPass:     args[2],
		SubscribeTopic: args[3],
		DBPath:         args[4],
	}
	if config.BrokerAddr == "" {
		return nil, fmt.Errorf("missing argument 0 broker_addr")
	}
	if config.BrokerUser == "" {
		return nil, fmt.Errorf("missing argument 1 broker_user")
	}
	if config.BrokerPass == "" {
		return nil, fmt.Errorf("missing argument 2 broker_pass")
	}
	if config.SubscribeTopic == "" {
		return nil, fmt.Errorf("missing argument 3 subscribe_topic")
	}
	if config.DBPath == "" {
		return nil, fmt.Errorf("missing argument 4 db_path")
	}
	log.Printf("[config] broker_addr: %q", config.BrokerAddr)
	log.Printf("[config] broker_user: %q", config.BrokerUser)
	log.Printf("[config] broker_pass: %q", config.BrokerPass)
	log.Printf("[config] subscribe_topic: %q", config.SubscribeTopic)
	log.Printf("[config] db_path: %q", config.DBPath)
	return config, nil
}
