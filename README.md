# home-temp-monitor

Pet project to monitor temperature and humidity over time across my house.

# Probe

The `/probe` folder contains the code for arduino-compatible devices that have a WiFi connection.

It works by using pin **17** for fetching environment data and publishes it to a mqtt broker with the following format :
```json
{
	"i": "room_id",
	"t": 123, // temperature : 12.3 °C
	"h": 456, // humidity :    45.6 %
	"d": 1234 // number of seconds elapsed since the measure
}
```

# Monitoring station

The `/station`  folder contains the monitoring server.
> It requires mosquitto to run on the same machine.

It first goal is to subscribe to the mqtt broker and persist measures.
Also, it provides an api and web interface to display the measures over time through charts.