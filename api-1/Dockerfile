FROM golang:latest

RUN adduser --disabled-password --gecos '' api-1
USER api-1

WORKDIR /go/src/app
COPY . .

RUN go get github.com/pilu/fresh
RUN go get ./...

CMD [ "fresh" ]