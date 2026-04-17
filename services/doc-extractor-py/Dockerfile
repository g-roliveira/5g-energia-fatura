# syntax=docker/dockerfile:1.7

FROM python:3.12-slim

ENV PYTHONDONTWRITEBYTECODE=1 \
    PYTHONUNBUFFERED=1 \
    DOC_EXTRACTOR_HOST=0.0.0.0 \
    DOC_EXTRACTOR_PORT=8090

WORKDIR /app

COPY pyproject.toml ./
COPY src ./src
COPY services/doc-extractor-py ./services/doc-extractor-py

RUN pip install --no-cache-dir --upgrade pip && \
    pip install --no-cache-dir ".[ocr]"

EXPOSE 8090

CMD ["python", "services/doc-extractor-py/app/main.py"]
