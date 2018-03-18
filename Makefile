
run:
	docker-compose up

test:
	docker-compose run kyuko-app go test -v ./...

build:
	docker build -t kyuko-bin -f Dockerfile.bin .
	docker run -v $(PWD)/bin:/go/src/github.com/g-hyoga/kyuko/bin kyuko-app

deploy: build
	docker build -t kyuko-app .
	docker tag kyuko-app:latest 878798127453.dkr.ecr.ap-northeast-1.amazonaws.com/kyuko-app:latest
	docker push 878798127453.dkr.ecr.ap-northeast-1.amazonaws.com/kyuko-app:latest

mac-build:
	docker build -t kyuko-bin -f Dockerfile.bin .
	docker run -v $(PWD)/bin:/go/src/github.com/g-hyoga/kyuko/bin -e GOOS=darwin -e GOARCH=amd64 kyuko-app
