# Build stage
FROM golang:1.22-alpine AS builder
WORKDIR /app

# h3-go requires CGO + C compiler
RUN apk add --no-cache git ca-certificates build-base

COPY go.mod ./
RUN go mod download

COPY . .
RUN go mod tidy

ENV CGO_ENABLED=1
RUN go build -ldflags="-s -w" -o /server ./cmd/server

# Runtime
FROM alpine:3.19
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /server /app/server
ENV PORT=8080
EXPOSE 8080
CMD ["/app/server"]
