FROM golang:1.24

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

RUN go install github.com/swaggo/swag/cmd/swag@latest

COPY . .

RUN /go/bin/swag init --generalInfo cmd/httpserver/main.go --output docs

RUN apt-get update && apt-get install -y postgresql-client curl

COPY scripts/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

RUN go build -o main ./cmd/httpserver

EXPOSE 8080

ENTRYPOINT ["/entrypoint.sh"]
CMD ["./main"]
