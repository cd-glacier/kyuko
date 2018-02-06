setup:
	mkdir log

run:
	docker-compose up

test:
	docker-compose run kyuko-app go test -v ./...

