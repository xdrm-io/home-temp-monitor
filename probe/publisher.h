#pragma once

#include "measure.h"
#include "wireless.h"
#include <PubSubClient.h>

#define PUB_BUF_SIZE 32

class Publisher {
	public:
		Publisher(Wireless w, const char* host, const char* user, const char* pass);

		void loop();

		void publish(const char* topic, const Measure& m);


	protected:
		void reconnect();

	private:
		const char*  m_user    {};
		const char*  m_pass    {};
		Wireless&    m_wireless;
		PubSubClient m_client  {};
		Buffer       m_buffer  { PUB_BUF_SIZE };
};