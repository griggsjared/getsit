FROM golang:1.23 AS builder

WORKDIR /app

RUN apt update && apt install -y npm && npm install -g tailwindcss

COPY . .

RUN npx tailwindcss -i ./web/main.css -o ./web/public/main.css -c ./web/tailwind.config.js --watch

RUN go mod download

RUN go install github.com/a-h/templ/cmd/templ@latest

RUN templ generate

RUN CGO_ENABLED=0 GOOS=linux go build -o /entrypoint ./cmd/web

FROM gcr.io/distroless/base-debian11 AS build

WORKDIR /

COPY --chown=nonroot --from=builder /entrypoint /app/web/public /

EXPOSE 8080

USER nonroot:nonroot

CMD ["/entrypoint"]
