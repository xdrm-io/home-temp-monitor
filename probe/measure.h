#pragma once

#include <Arduino.h>

struct Measure { // 8 bytes
	// (4 bytes) timestamp in seconds ; millis % 1e3
	unsigned long timestamp;
	// (2 bytes) temperature in degrees x10
	uint16_t      temperature;
	// (2 bytes) humidity in percents x10
	uint16_t      humidity;

	// SOLVE
	// function "Measure::operator=(const Measure &)" (declared implicitly) cannot be referenced -- it is a deleted functionC/C++(1776)
	void operator=(const Measure& m);
};

class Buffer {
	public:
		Buffer(const size_t capacity);
		~Buffer();

		void append(const Measure& m);
		void clear();

		const Measure* get(const size_t i) const;
		size_t         size() const;
		bool           full() const;

	private:
		Measure* m_data     { nullptr };
		size_t   m_size     { 0 };
		size_t   m_capacity { 0 };
};