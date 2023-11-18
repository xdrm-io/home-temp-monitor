# load .env
include .env
export $(sed 's/=.*//' .env)

.PHONY: build web run package

build:
	go build -ldflags '-s -w' -o station .

run:
	BROKER_ADDR="${BROKER_ADDR}" BROKER_USER="${BROKER_USER}" BROKER_PASS="${BROKER_PASS}" MQTT_TOPIC="${MQTT_TOPIC}" DB_PATH="${DB_PATH}" ./station

web:
	yarn --cwd endpoint/web && yarn --cwd endpoint/web build

package: web
	rm -rf ./pkg
	mkdir -p ./pkg

	GOOS=linux GOARCH=arm GOARM=5 CGO_ENABLED=0 go build -o ./pkg/station .
	chmod +x ./pkg/station

	echo '#!/bin/bash' > ./pkg/station.sh
	echo "root=\$$(dirname \$$(realpath \$$0))" >> ./pkg/station.sh
	echo "BROKER_ADDR=\"${BROKER_ADDR}\" BROKER_USER=\"${BROKER_USER}\" BROKER_PASS=\"${BROKER_PASS}\" MQTT_TOPIC=\"${MQTT_TOPIC}\" DB_PATH=\"\$${root}/${DB_PATH}\" \$${root}/station" >> ./pkg/station.sh
	chmod +x ./pkg/station.sh

	cp ./station.service ./pkg/station.service

	cp ./install.sh ./pkg/install.sh
	chmod +x ./pkg/station.sh

	touch ./pkg/db.sqlite
	chmod 777 ./pkg/db.sqlite
