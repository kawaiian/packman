FROM golang:alpine AS build

ENV GOPATH /go

COPY . /go/src/github.com/kawaiian/packman

WORKDIR /go/src/github.com/kawaiian/packman/cmd/packman

RUN go build -o packman packman.go

FROM alpine

COPY --from=build /go/src/github.com/kawaiian/packman/cmd/packman/packman /app/

WORKDIR /app

CMD ["./packman"]

