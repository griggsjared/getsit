FROM golang:1.23 AS builder

WORKDIR /app

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o ./getsit-api ./cmd/api

FROM gcr.io/distroless/base-debian11 AS build

COPY --from=builder /app/getsit-api /getsit-api

EXPOSE 8080

USER nonroot:nonroot

CMD ["/getsit-api"]
