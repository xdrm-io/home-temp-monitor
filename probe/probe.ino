// common
#define WIFI_SSID    "WIFI_SSID"
#define WIFI_PASS    "WIFI_PASS"
#define BROKER_ADDR  "192.168.0.xxx"
#define BROKER_USER  "mqtt_user"
#define BROKER_PASS  "mqtt_pass"
// device
#define ROOM_ID  "room_name"
#define BUF_SIZE 32


#include <DHT.h>
#include "wireless.h"
#include "publisher.h"

DHT sensor(17, DHT11);

Wireless  wireless(WIFI_SSID, WIFI_PASS);
Publisher publisher(wireless, BROKER_ADDR, BROKER_USER, BROKER_PASS);

void setup(){
	Serial.begin(115200);
	wireless.setup();
}

void loop() {
	// required for internal work
	publisher.loop();

	Measure m;
	m.timestamp = millis() % 1000;
	m.temperature = 10 * sensor.readTemperature(false, false);
	m.humidity = 10 * sensor.readHumidity(false);

	// no new data.
	if( m.temperature == NAN || m.humidity == NAN ){
		delay(1000);
		return;
	}
	publisher.publish(ROOM_ID, m);
}