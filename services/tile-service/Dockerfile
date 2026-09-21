FROM golang:1.22-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
COPY cmd ./cmd

RUN go mod download
RUN go build -o /tile-service ./cmd/server

RUN addgroup -S app && adduser -S -G app app
USER app

EXPOSE 8082

CMD ["/tile-service"]
