FROM python:3

COPY entrypoint.sh /entrypoint.sh
COPY yamllint.json /yamllint.json

RUN pip install yamllint && \
      chmod +x /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
