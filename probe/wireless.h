#pragma once

#include <ESP8266WiFi.h>

class Wireless {
	public:
		Wireless(const char* ssid, const char* pass);

		void setup();

		void reconnect();
		bool connected() const;

		WiFiClient* client();

	private:
		const char* m_ssid   { nullptr };
		const char* m_pass   { nullptr };
		WiFiClient* m_client { nullptr };
};