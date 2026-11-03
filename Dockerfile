FROM golang:1.26-alpine AS builder

WORKDIR /src
RUN apk add --no-cache ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/app ./cmd/app

FROM alpine:3.22

RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /out/app /app/app
COPY --from=builder /src/database/migrations /app/database/migrations

USER 65532:65532
EXPOSE 8080
ENTRYPOINT ["/app/app"]
