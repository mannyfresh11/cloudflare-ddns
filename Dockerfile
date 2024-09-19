FROM golang:1.23.1-alpine3.20

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download && go mod verify

COPY . .

RUN echo "Listing files after copying source code:" && ls -al /app

RUN go build cmd/main.go 

RUN echo "Listing files after copying source code:" && ls -al /app

CMD ["./main"]
