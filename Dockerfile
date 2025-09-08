FROM debian:latest


RUN mkdir /gorrent
COPY cmd go.mod internal config  go.sum /gorrent/
WORKDIR /gorrent

RUN mkdir /data
RUN mkdir /var/log/gorrent
RUN apt update
RUN apt-get update -y && apt-get install ca-certificates -y
RUN apt install golang-go -y
RUN go build cmd/superserver/main.go
RUN rm -rf cmd  config  go.mod  go.sum  internal

ENTRYPOINT ["./main"]

