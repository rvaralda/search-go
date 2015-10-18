FROM golang:1.5.1
ADD .  /app
WORKDIR /app
RUN go get -d
RUN go build -o es es.go
CMD ["/app/es"]
EXPOSE 9000
