FROM golang:1.22.5-alpine as builder

WORKDIR /build

RUN apk add --no-cache git make

RUN git clone https://github.com/StevenSermeus/goval.git

WORKDIR /build/goval

RUN git checkout 1.0.0

RUN make build

FROM alpine:3.14

ENV VERSION=1.0.0

WORKDIR /goval

COPY --from=builder /build/goval/goval /goval

ENTRYPOINT ["/goval/goval"]