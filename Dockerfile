FROM golang:1.9

WORKDIR /go/src/github.com/g-hyoga/kyuko
COPY . .

RUN go get -u github.com/golang/dep/...
RUN dep init 

ENV GOOS linux
ENV GOARCH amd64

CMD ["go", "build", "-o", "./bin/kyuko-lambda", "src/cmd/main.go"]
