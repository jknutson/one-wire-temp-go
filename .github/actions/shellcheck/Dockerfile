FROM ubuntu:21.04

COPY entrypoint.sh /entrypoint.sh
COPY shellcheck.json /shellcheck.json

# hadolint ignore=DL3008
RUN apt-get update && apt-get install -y --no-install-recommends shellcheck && \
      apt-get clean && \
      rm -rf /var/lib/apt/lists/*

ENTRYPOINT ["/entrypoint.sh"]

CMD ["."]
