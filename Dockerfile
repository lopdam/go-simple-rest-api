# syntax=docker/dockerfile:1
FROM golang:1.22.3

COPY go.mod .
COPY go.sum .

RUN go mod download

WORKDIR /code

COPY . /code/

CMD ["go", "run", "main.go"]