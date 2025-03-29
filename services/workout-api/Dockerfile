FROM golang:1.24.1-alpine

RUN apk add --no-cache bash curl

WORKDIR /app

COPY . .

RUN go build -o workout-api main.go 

EXPOSE 8080

CMD ["./workout-api"]
