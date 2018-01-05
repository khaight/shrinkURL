FROM golang:latest

WORKDIR /app

ENV SRC_DIR=/go/src/github.com/shrinkUrl/

ADD . $SRC_DIR

RUN cd $SRC_DIR; go build -o main; cp main /app/

ENTRYPOINT ["./main"]
