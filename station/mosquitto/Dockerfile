FROM eclipse-mosquitto:latest

ARG BROKER_USERNAME
ARG BROKER_PASSWORD
ARG PROBE_USERNAME
ARG PROBE_PASSWORD

ENV BROKER_USERNAME=${BROKER_USERNAME}
ENV BROKER_PASSWORD=${BROKER_PASSWORD}
ENV PROBE_USERNAME=${PROBE_USERNAME}
ENV PROBE_PASSWORD=${PROBE_PASSWORD}

COPY mosquitto.conf /mosquitto/config/mosquitto.conf

# generate passwords from env variables
RUN mosquitto_passwd -b -c /mosquitto/config/passwd "${BROKER_USERNAME}" "${BROKER_PASSWORD}" && \
	mosquitto_passwd -b /mosquitto/config/passwd "${PROBE_USERNAME}" "${PROBE_PASSWORD}"

EXPOSE 1883/tcp
CMD [ "/usr/sbin/mosquitto", "-c", "/mosquitto/config/mosquitto.conf" ]