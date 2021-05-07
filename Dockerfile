FROM golang:1.16

WORKDIR /app
COPY .  /app

RUN go build

CMD ["/app/search-go"]
EXPOSE 9000
