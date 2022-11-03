FROM golang:1.19

WORKDIR /go/app
COPY . .

RUN go mod download