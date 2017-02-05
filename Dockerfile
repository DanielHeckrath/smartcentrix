FROM golang:1.7.4-alpine

ENV  GOPATH /go
ENV APPPATH $GOPATH/src/github.com/DanielHeckrath/smartcentrix
COPY . $APPPATH
RUN cd $APPPATH && go build -o /smartcentrix-api && rm -rf $GOPATH

EXPOSE 8080
EXPOSE 8081
EXPOSE 8082

ENTRYPOINT ["/smartcentrix-api"]