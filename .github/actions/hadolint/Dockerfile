FROM hadolint/hadolint:latest

COPY entrypoint.sh /entrypoint.sh
COPY hadolint.json /hadolint.json

ENTRYPOINT ["/entrypoint.sh"]

CMD ["."]
