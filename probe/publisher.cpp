#include "publisher.h"

#include <string>

Publisher::Publisher(Wireless w, const char* host, const char* user, const char* pass)
: m_user(user), m_pass(pass), m_wireless(w), m_client(*w.client())
{
	m_client.setServer(host, 1883);
}

void Publisher::loop() {
	m_client.loop();
}

void Publisher::reconnect() {
	const auto& id = String(random(0xffff), HEX);

	while( !m_client.connected() ){
		Serial.println("[pub] reconnecting");
		if( !m_client.connect(id.c_str(), m_user, m_pass) ){
			Serial.println("[pub] failed to connect");
			delay(1000);
		}
	}
	Serial.println("[pub] connected");
}

void Publisher::publish(const char* room, const Measure& m){
	// no space remaining in buffer
	if( m_buffer.full() ){
		Serial.println("[pub] buffer already full");
		return;
	}
	m_buffer.append(m);

	if( !m_buffer.full() ){
		return;
	}

	// publish
	m_wireless.reconnect();
	reconnect();

	// range over buffer measures
	const auto now = millis();
	for( auto i = 0 ; i < m_buffer.size() ; i++ ){
		const auto measure = m_buffer.get(i);
		if( measure == nullptr ){
			Serial.println("[pub] measure is null");
			continue;
		}
		const auto& topic = String("/room/") + room + "/env";

		// json payload with:
		// - i: room id
		// - t: temperature
		// - h: humidity
		// - d: time diff in seconds from the measure event e.g. 3 = 3 seconds ago
		const auto& payload = String("{")
			+ "\"t\":" + String(measure->temperature) + ","
			+ "\"h\":" + String(measure->humidity) + ","
			+ "\"d\":" + String( (now-measure->timestamp) / 1000 )
			+ "}";

		Serial.print("[pub] topic '");
		Serial.print(topic);
		Serial.print("' payload '");
		Serial.print(payload);
		Serial.println("'");
		m_client.publish(topic.c_str(), payload.c_str());
	}
	m_buffer.clear();
}