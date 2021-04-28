FROM ubuntu:latest

COPY entrypoint.sh /entrypoint.sh
COPY goss.json /goss.json

RUN apt-get update && apt-get install -y curl && \
      curl -L https://github.com/aelsabbahy/goss/releases/latest/download/goss-linux-arm -o /usr/local/bin/goss && \
      # TODO: verify shasum
      # curl -L https://github.com/aelsabbahy/goss/releases/latest/download/goss-linux-arm.sha256 -o ./goss-linux-arm.sha256 && \
      # echo "$(cat ./goss-linux-arm.sha256) /usr/local/bin/goss" | sha256sum --check --status && \
      chmod +rx /usr/local/bin/goss && \
      chmod +rx /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]

CMD ["."]
