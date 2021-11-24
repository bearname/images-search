download:
	go mod download

build:
	go build -o face-server.exe cmd/main.go

run: build
	face-server.exe