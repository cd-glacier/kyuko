binary-name=kyukoHandler
lambda-name=kyuko
docker-image=kyuko-image

image-build:
	docker build -t $(docker-image) .

build: clean image-build
	docker run -v $(PWD)/bin:/go/src/github.com/g-hyoga/kyuko/bin kyuko-image 
	if [ ! -d output ]; then \
	  mkdir output; \
  fi
	cd bin && zip ../output/handler.zip $(binary-name)

local-build: clean
	GOOS=linux GOARCH=amd64 go build -o ./bin/$(binary-name) src/cmd/main.go 
	if [ ! -d output ]; then \
		mkdir output; \
	fi
	cd bin && zip ../output/handler.zip $(binary-name)

clean:
	rm -rf bin
	rm -rf output

# not working
test: 
	docker run $(docker-image) go test -v test ./...

local-test:
	go test -v ./...

# not working
deploy:
	aws cloudformation package \
		--template-file ./template.yml \
		--output-template-file output/output-template.yml \
		--s3-bucket kyuko-package
	aws cloudformation deploy \
		--template-file ./output/output-template.yml \
		--stack-name new-stack-name


