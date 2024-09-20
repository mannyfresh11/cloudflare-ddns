FROM golang:1.23.1-alpine3.20 AS builder

WORKDIR /tmp/cf

COPY go.mod .
COPY go.sum .

RUN go mod download && go mod verify

COPY . .

RUN go build cmd/cfupdater.go 

FROM alpine:3.20

COPY --from=builder /tmp/cf /app

CMD ["./app/cfupdater"]
