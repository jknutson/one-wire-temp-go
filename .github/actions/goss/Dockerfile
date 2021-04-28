FROM ubuntu:21.04

COPY entrypoint.sh /entrypoint.sh
COPY goss.json /goss.json

# hadolint ignore=DL3008
RUN apt-get update && apt-get install -y --no-install-recommends curl && \
      apt-get clean && \
      rm -rf /var/lib/apt/lists/* && \
      curl -L https://github.com/aelsabbahy/goss/releases/latest/download/goss-linux-arm -o /usr/local/bin/goss && \
      chmod +rx /usr/local/bin/goss && \
      chmod +rx /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]

CMD ["validate"]
