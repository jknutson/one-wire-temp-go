FROM hadolint/hadolint:latest

COPY entrypoint.sh /entrypoint.sh
COPY hadolint.json /hadolint.json

RUN chmod +x /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]

CMD ["."]
