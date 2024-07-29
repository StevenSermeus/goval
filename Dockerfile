FROM golang:1.22.5-alpine3.20 AS build

WORKDIR /compiler

COPY go.mod .

COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o goval ./goval.go

FROM scratch AS prod

WORKDIR /production

COPY --from=build /compiler/goval .

EXPOSE 8080

CMD ["./goval"]