FROM golang:1.14

LABEL maintainer="uday.manchanda14@gmail.com"

WORKDIR /go/src/app

COPY . /go/src

RUN go get -d -v ./...

RUN go install -v ./...

CMD ["go", "run", "src/app/main.go"]