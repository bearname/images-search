.PHONY: download migrateup migratedown lint build buildWin run up down herokulogs

download:
	go mod download

build:
	set GOARCH=amd64
	set GOOS=linux
	set CGO_ENABLED=0
	pwd
	go build -o bin/imageProcessor cmd/imageProcessor/main.go
	go build -o bin/liberatorBlockedWorker cmd/liberatorBlockedWorker/main.go
	go build -o bin/rawImageHandler cmd/rawImageHandler/main.go
	go build -o bin/photoservice cmd/backend/main.go
	ls

buildWin:
	set GOARCH=amd64
	set GOOS=windows
	set CGO_ENABLED=0
	go build -o imageProcessor.exe cmd/imageProcessor/main.go
	go build -o liberatorBlockedWorker.exe cmd/liberatorBlockedWorker/main.go
	go build -o rawImageHandler.exe cmd/rawImageHandler/main.go
	go build -o photoservice.exe cmd/backend/main.go

cleanFront:
	del /Q /S dist

buildFront:
	cd frontend && npm install && npm run build && cd .. && xcopy  /e "frontend\dist\*.*" "dist\*.*"

run: buildWin
	face-server.exe

migrateup:
	migrate -path ./data/migrations/pg -database "postgresql://url-short:1234@localhost:5432/url-short?sslmode=disable" -verbose up

migratedown:
	migrate -path ./data/migrations/pg -database "postgresql://url-short:1234@localhost:5432/url-short?sslmode=disable" -verbose down

lint:
	golangci-lint run

up: build
	docker-compose build --parallel
	docker-compose up -d

down:
	docker-compose down

reload:
	make down && make build

herokulogs:
	heroku logs --tail

profile:
	go tool pprof -inuse_space memtest http://localhost:8081/debug/pprof/heap

stripeTrigger:
	stripe trigger payment_intent.succeeded

stripeForward:
	stripe listen --forward-to http://localhost:8000/webhook

runNgrok:
	ngrok http -bind-tls=true localhost:8000

all: download buildFront build