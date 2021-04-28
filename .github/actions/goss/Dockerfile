FROM ubuntu:latest

RUN apt-get update && apt-get install -y curl && \
      curl -L https://github.com/aelsabbahy/goss/releases/latest/download/goss-linux-arm -o /usr/local/bin/goss && \
      curl -L https://github.com/aelsabbahy/goss/releases/latest/download/goss-linux-arm.sha256 -o ./goss-linux-arm.sha256 && \
      echo "$(cat ./goss-linux-arm.sha256) /usr/local/bin/goss" | sha256sum --check --status && \
      chmod +rx /usr/local/bin/goss

COPY entrypoint.sh /entrypoint.sh
COPY goss.json /goss.json

ENTRYPOINT ["/entrypoint.sh"]

CMD ["."]
