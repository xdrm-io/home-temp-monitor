# load .env
include .env
export $(sed 's/=.*//' .env)

.PHONY: build web run package

build:
	go build -ldflags '-s -w' -o station .

run:
	./station "${BROKER_ADDR}" "${BROKER_USER}" "${BROKER_PASS}" "${MQTT_TOPIC}" "${DB_PATH}"

web:
	yarn --cwd endpoint/web && yarn --cwd endpoint/web build

package: web
	rm -rf ./pkg
	mkdir -p ./pkg

	GOOS=linux GOARCH=arm GOARM=5 CGO_ENABLED=0 go build -o ./pkg/station .
	chmod +x ./pkg/station

	echo '#!/bin/bash' > ./pkg/station.sh
	echo "root=\$$(dirname \$$(realpath \$$0))" >> ./pkg/station.sh
	echo "\$${root}/station \"${BROKER_ADDR}\" \"${BROKER_USER}\" \"${BROKER_PASS}\" \"${MQTT_TOPIC}\" \$${root}/${DB_PATH}" >> ./pkg/station.sh
	chmod +x ./pkg/station.sh

	cp ./station.service ./pkg/station.service

	cp ./install.sh ./pkg/install.sh
	chmod +x ./pkg/station.sh

	touch ./pkg/db.sqlite
	chmod 777 ./pkg/db.sqlite
