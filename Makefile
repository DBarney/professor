unexport GOPATH
test:
	go test ./...

build: test
	@rm -rf ./bin
	@mkdir -p ./bin
	go build -o ./bin/prof ./*.go