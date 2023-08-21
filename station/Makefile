# load .env
include .env
export $(sed 's/=.*//' .env)

.PHONY: build run package

build:
	go build -ldflags '-s -w' -o station .

run:
	./station "${BROKER_ADDR}" "${BROKER_USER}" "${BROKER_PASS}" "${MQTT_TOPIC}"

package:
	rm -rf ./pkg
	mkdir -p ./pkg
	go build -ldflags '-s -w' -o ./pkg/station .
	# create station.sh that runs the binary with env vars with values from the current .env file
	echo '#!/bin/bash' > ./pkg/station.sh
	echo "./station \"${BROKER_ADDR}\" \"${BROKER_USER}\" \"${BROKER_PASS}\" \"${MQTT_TOPIC}\"" >> ./pkg/station.sh
	chmod +x ./pkg/station.sh

