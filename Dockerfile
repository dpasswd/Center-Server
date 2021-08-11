FROM golang:1.16.5

ENV GOPROXY https://goproxy.cn/

RUN mkdir -p /data/www/api-dpasswd

WORKDIR /data/www/api-dpasswd/

ADD . /data/www/api-dpasswd

RUN apt-get install gcc

RUN go build main.go

EXPOSE 8080

ENTRYPOINT ./main server
