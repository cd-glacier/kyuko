DBNAME:=kyuko
TESTDB:=kyuko_dev
ENV:=development

setup:
	go get github.com/tools/godep
	#go get github.com/rubenv/sql-migrate/...
	#go get github.com/PuerkitoBio/goquery
	#go get golang.org/x/text/encoding/japanese
	#go get golang.org/x/text/transform


build:
	go build -o cmd/main cmd/main.go

run:
	./cmd/main

test/init:
	mysql -u root -h localhost --protocol tcp -e "create database \`$(TESTDB)\`" -p

migrate/init:
	sudo mysql.server restart
	mysql -u root -h localhost --protocol tcp -e "create database \`$(DBNAME)\`" -p

migrate/up:
	sql-migrate up 

migrate/down:
	sql-migrate down

migrate/status:
	sql-migrate status
