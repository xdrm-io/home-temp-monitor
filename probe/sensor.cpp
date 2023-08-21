#include "sensor.h"

Sensor::Sensor(const uint8_t pin, const uint8_t type)
: m_sensor(pin, type){
}

Sensor::~Sensor(){
}

SensorData Sensor::read(){
	SensorData data;

	// x10 to make it an integer with 1 decimal
	data.temperature = 10 * m_sensor.readTemperature(false, false);
	data.humidity    = 10 * m_sensor.readHumidity(false);

	if( data.temperature == NAN || data.humidity == NAN ){
		return data;
	}

	data.ok = true;
	return data;
}