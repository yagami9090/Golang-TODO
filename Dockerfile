FROM golang:1.16-alpine

WORKDIR /app

COPY  . /app/

RUN go mod init todo && go mod tidy

CMD [ "go" , "run" , "main.go" ]
