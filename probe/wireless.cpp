#include "wireless.h"

Wireless::Wireless(const char* ssid, const char* pass)
: m_ssid(ssid), m_pass(pass), m_client(new WiFiClient())
{}

void Wireless::setup(){
	WiFi.mode(WIFI_STA);
	WiFi.begin(
		const_cast<char*>(m_ssid),
		const_cast<char*>(m_pass)
	);
}

bool Wireless::connected() const {
	return (WiFi.status() == WL_CONNECTED);
}

void Wireless::reconnect() {
	if( connected() ){
		return;
	}

	Serial.print("[wifi] connecting to ");
	Serial.print(m_ssid);
	while( !connected() ){
			delay(500);
			Serial.print(".");
	}
	Serial.println();

	Serial.println("[wifi] connected");
	Serial.println("[wifi] IP: ");
	Serial.println(WiFi.localIP());
}

WiFiClient* Wireless::client() {
	return m_client;
}