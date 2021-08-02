FROM golang:1.16-alpine

ADD . /go/src/arigo/bank/app

WORKDIR /go/src/arigo/bank/app

COPY go.mod ./

COPY go.sum ./

RUN go mod download

COPY *.go ./

RUN GOOS=linux go build -o /go/src/arigo/bank/app/bank-account /go/src/arigo/bank/app/main.go

EXPOSE 8300 

CMD ["/go/src/arigo/bank/app/bank-account"]