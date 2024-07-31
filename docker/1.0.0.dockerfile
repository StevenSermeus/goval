FROM golang:1.22.5-alpine as builder

WORKDIR /build

RUN apk add --no-cache git make

ENV COMMIT_SHA=a9bcded9d8245573086c5073beea540fd0d3d827

RUN git clone https://github.com/StevenSermeus/goval.git

WORKDIR /build/goval

RUN git checkout $COMMIT_SHA

RUN make build

FROM alpine:3.14

WORKDIR /goval

COPY --from=builder /build/goval/goval /goval

ENTRYPOINT ["/goval/goval"]