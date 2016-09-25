DBNAME:=kyuko
ENV:=development

setup:
	go get github.com/rubenv/sql-migrate/...

build:
	go build -o cmd/mai cmd/main.go

migrate/init:
	sudo mysql.server restrat
	mysql -u root -h localhost --protocol tcp -e "create database \`$(DBNAME)\`" -p

migrate/up:
	sql-migrate up 

migrate/down:
	sql-migrate down

migrate/status:
	sql-migrate status
