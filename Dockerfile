FROM golang:1.24-alpine AS base
FROM base AS build
WORKDIR /app

COPY go.mod go.sum ./
COPY cmd ./cmd
COPY internal ./internal

RUN CGO_ENABLED=0 go build -o cert-sync cmd/main/main.go

FROM base AS runner
WORKDIR /
COPY --from=build /app/cert-sync .
ENTRYPOINT ["./cert-sync"]
