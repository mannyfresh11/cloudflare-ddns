FROM golang:1.23.1-alpine3.20

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download && go mod verify

COPY . .

RUN go build cmd/cfupdater.go 

CMD ["./cfupdater"]
