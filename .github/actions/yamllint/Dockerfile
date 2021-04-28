FROM python:3

COPY entrypoint.sh /entrypoint.sh
COPY yamllint.json /yamllint.json

# hadolint ignore=DL3013
RUN pip install --no-cache-dir yamllint && \
      chmod +x /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]

CMD ["."]
