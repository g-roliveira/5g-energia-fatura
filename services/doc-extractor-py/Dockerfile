# syntax=docker/dockerfile:1.7

FROM python:3.12-slim

ENV PYTHONDONTWRITEBYTECODE=1 \
    PYTHONUNBUFFERED=1 \
    DOC_EXTRACTOR_HOST=0.0.0.0 \
    DOC_EXTRACTOR_PORT=8090

WORKDIR /app

# Copy root package (fatura module from src/) and its pyproject.toml
COPY pyproject.toml ./
COPY src ./src

# Install root package with OCR extras (includes mistralai)
RUN pip install --no-cache-dir --upgrade pip && \
    pip install --no-cache-dir ".[ocr]"

# Copy and install doc-extractor-py service package
COPY services/doc-extractor-py ./service
RUN pip install --no-cache-dir "./service"

EXPOSE 8090

CMD ["python", "-m", "app.main"]
