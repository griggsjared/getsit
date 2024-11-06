FROM golang:1.23 AS builder

WORKDIR /app

COPY . .

RUN wget https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-linux-x64 -O tailwindcss

RUN chmod +x tailwindcss

RUN ./tailwindcss -i /app/web/main.css -o /app/web/public/main.css -c /app/web/tailwind.config.js --minify

RUN go mod download

RUN go install github.com/a-h/templ/cmd/templ@latest

RUN templ generate

RUN CGO_ENABLED=0 GOOS=linux go build -o ./getsit-web ./cmd/web

FROM gcr.io/distroless/base-debian11 AS build

COPY --from=builder /app/getsit-web /getsit-web

COPY --from=builder /app/web/public /web/public

EXPOSE 8080

USER nonroot:nonroot

CMD ["/getsit-web"]
