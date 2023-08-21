// common
#define WIFI_SSID    "WIFI_SSID"
#define WIFI_PASS    "WIFI_PASS"
#define BROKER_ADDR  "192.168.1.10"
#define BROKER_USER  "probe"
#define BROKER_PASS  "probe_pass"
// device
#define ROOM_ID  "room_name"
#define BUF_SIZE 32


#include "sensor.h"
#include "wireless.h"
#include "publisher.h"

Sensor sensor(17, DHT11);

Wireless  wireless(WIFI_SSID, WIFI_PASS);
Publisher publisher(wireless, BROKER_ADDR, BROKER_USER, BROKER_PASS);

void setup(){
	randomSeed(micros());

	Serial.begin(115200);
	Serial.println();
	wireless.setup();
	wireless.reconnect();
}

void loop() {
	// required for internal work
	publisher.loop();

	const auto& data = sensor.read();
	if( !data.ok ){  // no new data
		delay(1000);
		return;
	}
	Measure m;
	m.timestamp = millis() % 1000;
	m.temperature = data.temperature;
	m.humidity = data.humidity;

	Serial.print("[sensor] ");
	Serial.print(m.temperature);
	Serial.print(" / ");
	Serial.println(m.humidity);
	publisher.publish(ROOM_ID, m);
}