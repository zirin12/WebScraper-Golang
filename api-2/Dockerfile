FROM golang:latest

RUN adduser --disabled-password --gecos '' api-2
USER api-2

WORKDIR /go/src/app
COPY . .

RUN go get github.com/pilu/fresh
RUN go get ./...

CMD [ "fresh" ]