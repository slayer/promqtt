# MQTT To Prometheus bridge


[![Build Status](https://travis-ci.org/slayer/promqtt.png?branch=master)](https://travis-ci.org/slayer/promqtt)
[![Go Report Card](https://goreportcard.com/badge/github.com/slayer/promqtt)](https://goreportcard.com/report/github.com/slayer/promqtt)
[![Docker Pulls](https://img.shields.io/docker/pulls/slayer/promqtt.svg?maxAge=604800)](https://hub.docker.com/r/slayer/promqtt/)

Send MQTT data to Prometheus.

Based on https://github.com/slayer/mqtt_blackbox_exporter

## Build

```sh
$ mkdir -p ${GOPATH}/src/github.com/slayer/
$ git clone https://github.com/slayer/promqtt.git ${GOPATH}/src/github.com/slayer/promqtt/
$ cd ${GOPATH}/src/github.com/slayer/promqtt/
$ make
```

This will build the promqtt for all target platforms and write them to the ``build/`` directory.

TODO: upload binaries
Binaries are provided on Github, see https://github.com/slayer/promqtt.

## Install

Place the binary somewhere in a ``PATH`` directory and make it executable (``chmod +x promqtt``).

## Configure

See ``config.yaml.dist`` for a configuration example.

## Run

Native:

```sh
$ ./promqtt -config.file config.yaml
```

Using Docker:

```sh
docker run --rm -it -p 9214:9214 -v ${PWD}/:/data/ slayer/promqtt:latest -config.file /data/config.yaml
```

```
$ curl -s http://127.0.0.1:9214/metrics

# TODO: paste output
```
