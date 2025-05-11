FROM --platform=$BUILDPLATFORM golang:1.24-alpine AS base
ARG BUILDPLATFORM

FROM base AS build
ARG TARGETARCH
ARG TARGETOS
ARG TARGETVARIANT
WORKDIR /app

COPY go.mod go.sum ./
COPY cmd ./cmd
COPY internal ./internal

RUN CGO_ENABLED=0 GOOS="${TARGETOS}" GOARCH="${TARGETARCH}" GOARM="${TARGETVARIANT//v}" go build -o cert-sync cmd/main/main.go

FROM scratch
COPY --from=alpine:latest /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
WORKDIR /
COPY --from=build /app/cert-sync .
ENTRYPOINT ["./cert-sync"]
