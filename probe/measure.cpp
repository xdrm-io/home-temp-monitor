#include "measure.h"

void Measure::operator=(const Measure& m) {
	timestamp   = m.timestamp;
	temperature = m.temperature;
	humidity    = m.humidity;
}

Buffer::Buffer(const size_t capacity)
: m_capacity(capacity)
{
	m_data = reinterpret_cast<Measure*>( malloc(m_capacity * sizeof(Measure)) );
}

Buffer::~Buffer() {
	free(m_data);
}

void Buffer::append(const Measure& m) {
	if( m_size == m_capacity ){
		return;
	}

	m_data[m_size++] = m;
}

void Buffer::clear() {
	m_size = 0;
}

const Measure* Buffer::get(const size_t i) const {
	if( i >= m_size ){
		return nullptr;
	}
	return &m_data[i];
}

size_t Buffer::size() const {
	return m_size;
}

bool Buffer::full() const {
	return (m_size >= m_capacity);
}