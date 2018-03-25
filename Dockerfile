FROM golang:1.8 

COPY ./bin /usr/local/bin

CMD ["/usr/local/bin/main"]
