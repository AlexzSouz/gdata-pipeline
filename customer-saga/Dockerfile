# BUILD STAGE
FROM golang:1.20.4-alpine3.17 AS build
WORKDIR /go/src/saga

COPY ./customer-saga .

# RUNTIME STAGE
FROM build AS runtime
WORKDIR /go/src/saga

RUN addgroup -g 1111 bots
RUN adduser -u 1001 -G bots -D gbot
RUN chown -hR 1001:1111 /go/src/saga

ENTRYPOINT ["go", "run", "program.go"]