binary-name=kyukoHandler
lambda-name=kyuko
docker-image=kyuko-image
output=$(PWD)/output
zip-name=handler.zip

init: image-build

image-build:
	docker build -t $(docker-image) .

build: clean
	docker run -v $(PWD)/bin:/go/src/github.com/g-hyoga/kyuko/bin kyuko-image 
	if [ ! -d output ]; then \
	  mkdir output; \
  fi
	cd bin && zip $(output)/$(zip-name) $(binary-name)

local-build: clean
	GOOS=linux GOARCH=amd64 go build -o ./bin/$(binary-name) src/cmd/main.go 
	if [ ! -d output ]; then \
		mkdir output; \
	fi
	cd bin && zip $(output)/$(zip-name) $(binary-name)

clean:
	rm -rf bin
	rm -rf output

# not working
test: 
	docker run $(docker-image) go test -v test ./...

local-test:
	go test -v ./...

deploy:
	aws lambda update-function-code \
		--function-name $(lambda-name) \
		--zip-file fileb://$(output)/$(zip-name)

invoke:
	aws lambda invoke \
		--function-name $(lambda-name) \
		$(output)/output.json

