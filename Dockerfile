FROM debian:jessie
RUN apt-get update && apt-get install -y ca-certificates
COPY promqtt /bin/promqtt
ENTRYPOINT ["/bin/promqtt"]
CMD ["-config.file /config.yaml"]
