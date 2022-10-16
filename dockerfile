FROM golang:1.19.2-alpine3.16

RUN apk update && apk add --no-cache git && apk add --no-cach bash && apk add build-base

RUN mkdir /app
WORKDIR /app

COPY . .

RUN go get -d -v ./...

RUN go install -v ./...

RUN go build -o /build cmd/leaderboard/main.go

EXPOSE 8080

CMD ["/build"]