#pragma once

#include "measure.h"
#include "wireless.h"
#include <PubSubClient.h>

#define PUB_BUF_SIZE 32

class Publisher {
	public:
		Publisher(Wireless w, const char* host, const char* user, const char* pass, const char* pub_topic);

		void loop();

		void publish(const char* topic, const Measure& m);


	protected:
		void reconnect();

	private:
		const char*  m_user     {};
		const char*  m_pass     {};
		const char*  m_topic    {};
		Wireless&    m_wireless;
		PubSubClient m_client   {};
		Buffer       m_buffer   { PUB_BUF_SIZE };
};