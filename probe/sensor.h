#pragma once

#include <DHT.h>

struct SensorData {
	bool     ok { false };
	uint16_t temperature;
	uint16_t humidity;
};

class Sensor {
	public:
		Sensor(const uint8_t pin, const uint8_t type);
		~Sensor();

		SensorData read();

	private:
		DHT m_sensor;
};