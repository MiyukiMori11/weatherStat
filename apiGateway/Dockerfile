FROM golang:1.19.4-alpine3.16 as build

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -o /usr/local/bin/app 

CMD ["app", "run"]