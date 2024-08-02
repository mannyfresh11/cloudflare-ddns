FROM golang:1.22.5-alpine3.20

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod veriy

COPY . .
RUN go build -v -o /usr/local/bin/app ./...
