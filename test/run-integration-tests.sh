#!/bin/bash

set -e

service_pid=

cleanup() {
	rc=$?
	rm -f test/out.log
	if [ ! -z "${service_pid}" ]; then
		kill $service_pid
	fi
	exit $rc
}
trap cleanup INT TERM

echo "=> Starting exporter"
./promqtt -config.file config.yaml.dist &
service_pid=$!

echo "=> Waiting 5s"
sleep 5

echo "=> Requesting /metrics"
curl --silent --max-time 2 http://localhost:9214/metrics > test/out.log

echo "=> Killing exporter (pid=${service_pid})"
kill $service_pid

echo "=> Checking result"
grep 'promqtt_mqtt_messages_received_total [[:digit:]]' test/out.log
