FROM golang:1.8

WORKDIR /go/src/github.com/g-hyoga/kyuko
COPY . .

WORKDIR /go/src/github.com/g-hyoga/kyuko/src

CMD ["go", "run", "cmd/main.go"]
