# Render builds from repo root by default — this Dockerfile copies backend/
FROM golang:1.22-alpine AS builder
WORKDIR /app

# h3-go requires CGO + a C compiler (gcc via build-base)
RUN apk add --no-cache git ca-certificates build-base

COPY backend/go.mod ./
RUN go mod download

COPY backend/ .
RUN go mod tidy

ENV CGO_ENABLED=1
RUN go build -ldflags="-s -w" -o /server ./cmd/server

FROM alpine:3.19
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /server /app/server
ENV PORT=8080
EXPOSE 8080
CMD ["/app/server"]
