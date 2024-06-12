FROM golang:1.20-alphin

WORKDIR /app

COPY . .

RUN go build -o main .

CMD ["./cmd"]