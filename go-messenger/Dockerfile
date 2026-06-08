FROM golang:1.26-alpine

RUN apk add --no-cache git

WORKDIR /app

COPY . /app/

RUN mkdir -p bin && \
    go build -o bin/mgclient ./mongodb-client/client && \
    go build -o bin/pgclient ./postgres-client/client

CMD [ "/app/bin/mgclient" ]