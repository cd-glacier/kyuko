binary-name=kyuko-lambda
lambda-name=kyuko

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
	aws cloudformation package \
		--template-file ./template.yml \
		--output-template-file output/output-template.yml \
		--s3-bucket kyuko-package
	aws cloudformation deploy \
		--template-file ./output/output-template.yml \
		--stack-name new-stack-name


