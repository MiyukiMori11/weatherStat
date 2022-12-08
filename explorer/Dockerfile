FROM golang:1.19.4-alpine3.16 as build

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -o /usr/local/bin/app 

FROM build as explorer

CMD ["app", "run"]

FROM build as migration

RUN apk add --no-cache bash bind-tools gettext \
    && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

ENV cntDayArchive=5

ENTRYPOINT ["./entrypoint.sh"]