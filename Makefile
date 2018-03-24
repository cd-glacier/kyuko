binary-name=kyukoHandler
lambda-name=kyuko

build: clean
	docker build -t kyuko-image .
	docker run -v $(PWD)/bin:/go/src/github.com/g-hyoga/kyuko/bin kyuko-image 
	if [ ! -d output ]; then \
		mkdir output; \
	fi
	cd bin && zip ../output/handler.zip $(binary-name)

local-buid: clean
	go build -o ./bin/$(binary-name) src/cmd/main.go 
	if [ ! -d output ]; then \
		mkdir output; \
	fi
	cd bin && zip ../output/handler.zip $(binary-name)

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


