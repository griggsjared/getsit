FROM golang:1.23 AS builder

WORKDIR /app

COPY . .

RUN wget https://github.com/tailwindlabs/tailwindcss/releases/download/v4.0.0/tailwindcss-linux-x64 -O tailwindcss

RUN chmod +x tailwindcss

RUN ./tailwindcss -i /app/web/main.css -o /app/web/public/main.css --minify

RUN go mod download

RUN go install github.com/a-h/templ/cmd/templ@v0.3.819

RUN templ generate

RUN CGO_ENABLED=0 GOOS=linux go build -o ./getsit-web ./cmd/web

FROM gcr.io/distroless/base-debian11 AS build

COPY --from=builder /app/getsit-web /getsit-web

EXPOSE 8080

USER nonroot:nonroot

CMD ["/getsit-web"]
