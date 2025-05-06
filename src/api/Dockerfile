# 1. build stage
FROM golang:1.23-alpine AS build
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o lead-locate-api ./cmd/main.go

# 2. final stage
FROM alpine:3.12
WORKDIR /app
COPY --from=build /app/lead-locate-api /app/
CMD ["./lead-locate-api"]