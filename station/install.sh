#!/bin/bash
if [ "$(pwd)" != "/opt/station" ]; then
	echo "must be located in /opt/station/";
	exit 1;
fi;

sudo ln -fs /opt/station/station.service /etc/systemd/system/station.service
sudo systemctl enable station.service
sudo systemctl start station.service
sudo systemctl status station.service

