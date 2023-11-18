package main

import (
	"fmt"
	"log"
	"os"
)

type Config struct {
	BrokerAddr     string
	BrokerUser     string
	BrokerPass     string
	SubscribeTopic string
	DBPath         string
}

func ReadConfig() (*Config, error) {
	config := &Config{
		BrokerAddr:     os.Getenv("BROKER_ADDR"),
		BrokerUser:     os.Getenv("BROKER_USER"),
		BrokerPass:     os.Getenv("BROKER_PASS"),
		SubscribeTopic: os.Getenv("MQTT_TOPIC"),
		DBPath:         os.Getenv("DB_PATH"),
	}
	if config.BrokerAddr == "" {
		return nil, fmt.Errorf("missing env variable 'BROKER_ADDR'")
	}
	if config.BrokerUser == "" {
		return nil, fmt.Errorf("missing env variable 'BROKER_USER'")
	}
	if config.BrokerPass == "" {
		return nil, fmt.Errorf("missing env variable 'BROKER_PASS'")
	}
	if config.SubscribeTopic == "" {
		return nil, fmt.Errorf("missing env variable 'MQTT_TOPIC'")
	}
	if config.DBPath == "" {
		return nil, fmt.Errorf("missing env variable 'DB_PATH'")
	}
	log.Printf("[config] broker_addr: %q", config.BrokerAddr)
	log.Printf("[config] broker_user: %q", config.BrokerUser)
	log.Printf("[config] broker_pass: %q", config.BrokerPass)
	log.Printf("[config] subscribe_topic: %q", config.SubscribeTopic)
	log.Printf("[config] db_path: %q", config.DBPath)
	return config, nil
}
