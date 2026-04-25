# syntax=docker/dockerfile:1.7

FROM --platform=$BUILDPLATFORM golang:1.25-bookworm AS go-builder

ARG TARGETOS
ARG TARGETARCH

WORKDIR /src
COPY services/backend-go/go.mod services/backend-go/go.sum ./services/backend-go/
WORKDIR /src/services/backend-go
RUN go mod download
WORKDIR /src
COPY services/backend-go ./services/backend-go
WORKDIR /src/services/backend-go
RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} go build -o /out/backend-go ./cmd/api

FROM mcr.microsoft.com/playwright/python:v1.55.0-noble

ENV PYTHONDONTWRITEBYTECODE=1 \
    PYTHONUNBUFFERED=1 \
    BACKEND_HOST=0.0.0.0 \
    BACKEND_PORT=8080 \
    EXTRACTOR_BASE_URL=http://doc-extractor:8090 \
    BACKEND_INTEGRATION_PG_URL= \
    ARTIFACTS_DIR=/data/artifacts \
    BOOTSTRAP_PYTHON_BIN=python \
    BOOTSTRAP_SCRIPT_PATH=/app/scripts/bootstrap_neoenergia_token.py

WORKDIR /app

COPY pyproject.toml ./
COPY src ./src
COPY scripts ./scripts
RUN pip install --no-cache-dir --upgrade pip && \
    pip install --no-cache-dir ".[ocr]"

COPY --from=go-builder /out/backend-go /usr/local/bin/backend-go

RUN mkdir -p /data /app/artifacts

EXPOSE 8080

CMD ["backend-go"]
