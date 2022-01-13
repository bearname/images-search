.PHONY: download migrateup migratedown lint build buildWin run up down herokulogs

download:
	go mod download

build:
	set GOARCH=amd64
	set GOOS=linux
	set CGO_ENABLED=0
	go build -o bin/photoservice cmd/main.go

buildWin:
	set GOARCH=amd64
	set GOOS=windows
	set CGO_ENABLED=0
	go build -o face-server.exe cmd/main.go

cleanFront:
	del /Q /S dist

buildFront:
	cd frontend && npm install && npm run build && cd .. && xcopy  /e "frontend\dist\*.*" "dist\*.*"

run: buildWin
	face-server.exe

migrateup:
	migrate -path ./data/migrations/pg -database "postgresql://url-short:1234@localhost:5432/test?sslmode=disable" -verbose up

migratedown:
	migrate -path ./data/migrations/pg -database "postgresql://url-short:1234@localhost:5432/test?sslmode=disable" -verbose down

lint:
	golangci-lint run

up: build
	docker-compose up -d --build

down:
	docker-compose down

herokulogs:
	heroku logs --tail

all: download buildFront build
