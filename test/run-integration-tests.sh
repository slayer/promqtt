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
./promqtt -config.file config.yaml.dist -log.level=debug &
service_pid=$!
mosquitto_pub -t user@domain.com/001122334455/dhtt1 -r -m 12.34

echo "=> Waiting 1s"
sleep 1

echo "=> Requesting /metrics"
curl --silent --max-time 2 http://localhost:9214/metrics > test/out.log

echo "=> Killing exporter (pid=${service_pid})"
kill $service_pid

echo "=> Checking result"
grep 'promqtt_dhtt1{device="001122334455",user="user@domain.com"} 12.34' test/out.log
