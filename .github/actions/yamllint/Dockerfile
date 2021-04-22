FROM python:3

COPY entrypoint.sh /entrypoint.sh

RUN pip install yamllint && \
      chmod +x /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
