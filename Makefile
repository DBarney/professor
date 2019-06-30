test:
	go mod download
	go mod verify
	go test ./...

build: test
	@rm -rf ./bin
	@mkdir -p ./bin
	go build -o ./bin/prof ./*.go

image:test
	env GOOS=linux go build -o ./bin/prof ./*.go
	docker build -t dbarney/professor .