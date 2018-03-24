binary-name=kyuko-lambda

build:
	docker build -t kyuko-image .
	docker run -v $(PWD)/bin:/go/src/github.com/g-hyoga/kyuko/bin kyuko-image 
	mkdir output
	zip ./output/handler.zip ./bin/$(binary-name)

local-buid:
	go build -o ./bin/$(binary-name) src/cmd/main.go 
	mkdir output
	zip ./output/handler.zip ./bin/$(binary-name)

clean:
	rm -rf bin
	rm -rf output

deploy:


