
# Build Nuxt
FROM node:20-alpine AS frontend-builder
WORKDIR  /app
RUN corepack enable && corepack prepare pnpm@8.15.9 --activate
COPY frontend/package.json frontend/pnpm-lock.yaml ./
RUN pnpm install --frozen-lockfile --shamefully-hoist
COPY frontend .
RUN pnpm build

# Build API
FROM golang:1.24-alpine AS builder
ARG BUILD_TIME
ARG COMMIT
ARG VERSION
RUN apk update && \
    apk upgrade && \
    apk add --update git build-base gcc g++

WORKDIR /go/src/app
COPY ./backend/go.mod ./backend/go.sum ./
RUN go mod download
COPY ./backend .
RUN rm -rf ./app/api/public
COPY --from=frontend-builder /app/.output/public ./app/api/static/public
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags "-s -w -X main.commit=$COMMIT -X main.buildTime=$BUILD_TIME -X main.version=$VERSION"  \
    -o /go/bin/api \
    -v ./app/api/*.go

# Production Stage
FROM alpine:latest

ENV CATALOG_MODE=production
ENV CATALOG_STORAGE_DATA=/data/
ENV CATALOG_STORAGE_SQLITE_URL=/data/catalog.db?_pragma=busy_timeout=2000&_pragma=journal_mode=WAL&_fk=1

RUN apk --no-cache add ca-certificates
RUN mkdir /app
COPY --from=builder /go/bin/api /app

RUN chmod +x /app/api

LABEL Name=catalog Version=0.0.1
LABEL org.opencontainers.image.source="https://github.com/sdmad737/catalog"
EXPOSE 7745
WORKDIR /app
VOLUME [ "/data" ]

ENTRYPOINT [ "/app/api" ]
CMD [ "/data/config.yml" ]
