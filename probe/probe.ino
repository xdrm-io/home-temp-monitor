// common
#define WIFI_SSID    "WIFI_SSID"
#define WIFI_PASS    "WIFI_PASS"
#define BROKER_ADDR  "192.168.1.10"
#define BROKER_USER  "probe"
#define BROKER_PASS  "probe_pass"
// device
#define READ_PERIOD_MS 60*1000

#include "sensor.h"
#include "wireless.h"
#include "publisher.h"

Sensor sensor(D4, DHT22);

Wireless  wireless(WIFI_SSID, WIFI_PASS);
Publisher publisher(wireless, BROKER_ADDR, BROKER_USER, BROKER_PASS);

void setup(){
	randomSeed(micros());

	Serial.begin(115200);
	Serial.println();
	wireless.setup();
	wireless.reconnect();
	sensor.setup();
}

unsigned long last = millis();

void loop() {
	// required for internal work
	publisher.loop();

	const auto now = millis();
	if( now - last  < READ_PERIOD_MS ){
		delay(100);
		return;
	}
	last = now;

	const auto& data = sensor.read();
	if( !data.ok ){  // no new data
		return;
	}
	Measure m;
	m.timestamp = millis();
	m.temperature = data.temperature;
	m.humidity = data.humidity;

	Serial.print("[sensor] ");
	Serial.print(m.temperature/10.);
	Serial.print(" / ");
	Serial.println(m.humidity/10.);
	publisher.publish(ROOM_ID, m);
}