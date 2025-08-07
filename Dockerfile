# -------------- builder container --------------
FROM golang:1.22 AS builder

WORKDIR /go/src/

ARG VERSION

COPY go.mod .
COPY go.sum .

RUN go env -w GO111MODULE=on
RUN go env -w GOPROXY=https://mirrors.cloud.tencent.com/go/,direct

RUN go mod download

COPY . .

RUN make build

# -------------- runner container --------------
FROM alpine:3.20 AS runner

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tencent.com/g' /etc/apk/repositories

RUN apk --update --no-cache add bash

WORKDIR /app

COPY --from=builder /go/src/blueapps-go /usr/bin/blueapps-go

COPY --from=builder /go/src/templates /app/templates

# docs
RUN mkdir -p /app/apidocs

COPY --from=builder /go/src/pkg/docs /app/apidocs

ENV API_DOC_FILE_BASE_DIR=/app/apidocs

# templates
ENV TMPL_FILE_BASE_DIR=/app/templates

COPY --from=builder /go/src/static /app/static

# i18n
ENV I18N_FILE_BASE_DIR=/app/i18n

COPY --from=builder /go/src/i18n /app/i18n

# static files
ENV STATIC_FILE_BASE_DIR=/app/static

# logs
RUN mkdir -p /app/v3logs

ENV LOG_BASE_DIR=/app/v3logs

