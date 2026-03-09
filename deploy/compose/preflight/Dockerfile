FROM python:3.12-alpine

WORKDIR /app

COPY validate_config.py /app/validate_config.py

RUN addgroup -S app && adduser -S -G app app

USER app

ENTRYPOINT ["python", "/app/validate_config.py"]
