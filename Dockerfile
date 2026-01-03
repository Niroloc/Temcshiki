FROM golang:1.25

COPY . /app

WORKDIR /app

CMD ["go", "run", "src/main.go"]