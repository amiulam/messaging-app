FROM golang:1.23-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod tidy

COPY . .

RUN go build -o messaging-app

RUN chmod +x messaging-app

EXPOSE 4000

EXPOSE 8080

CMD ["./messaging-app"]