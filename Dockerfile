FROM debian:latest


RUN mkdir /gorrent
WORKDIR /gorrent
COPY . .

RUN mkdir /data
RUN mkdir /var/log/gorrent
RUN apt update
RUN apt-get update -y && apt-get install ca-certificates -y
RUN apt install golang-go -y
RUN go build cmd/superserver/main.go
RUN rm -r cmd/  config/  data/  Dockerfile  go.mod  gorrent_conf.json  gorrent.log  go.sum  internal/  LICENSE  README.md  server.log  test/ .git/ .env .gitignore .github/

ENTRYPOINT ["./main"]

