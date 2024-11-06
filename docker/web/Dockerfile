FROM golang:1.23 AS builder

WORKDIR /app

RUN apt update && apt install -y npm && npm install -g tailwindcss

COPY . .

RUN npx tailwindcss -i ./web/main.css -o ./web/public/main.css -c ./web/tailwind.config.js --watch

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o ./getsit-web ./cmd/web

FROM gcr.io/distroless/base-debian11 AS build

EXPOSE 8080

COPY --from=builder /app/getsit-web /app/web/public /

USER nonroot:nonroot

CMD ["/getsit-web"]
