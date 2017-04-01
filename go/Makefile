DBNAME:=kyuko
TESTDB:=kyuko_dev
ENV:=development

setup:
	go get github.com/rubenv/sql-migrate/...
	go get github.com/Sirupsen/logrus
	go get github.com/tools/godep
	godep get

build:
	godep go build -o cmd/main cmd/main.go

run:
	./cmd/main

test:
	godep go test -v ./... | grep --color=auto -A 3 -e FAIL || echo "PASS ALL TEST"

migrate/init:
	#sudo service mysql restart
	sudo mysql.server restart
	mysql -u root -h localhost --protocol tcp -e "create database \`$(DBNAME)\`" -p

migrate/up:
	sql-migrate up -env="production"

migrate/down:
	sql-migrate down -env="production"

migrate/status:
	sql-migrate status
