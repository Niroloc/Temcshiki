FROM golang:1.18

ADD src /app

WORKDIR /app

RUN go build -v -o main.go

CMD ["main"]