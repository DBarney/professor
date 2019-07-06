# super simplistic approach to detecting changes
SRC_FILES := $(shell find . -path ./.git -prune -o -name '*.go' -print)

# actual targets meant to be run
test: .test
image: .docker
build: bin/prof

.test: www/www.go $(SRC_FILES)
	go mod download
	go mod verify
	go vet ./...
	go test ./...
	@touch .test

www/www.go: www/index.html
	@echo rebuilding static assets
	@go-bindata -fs -pkg www '-ignore=.*[.]go' -prefix www/ -o ./www/www.go ./www/...

bin/prof: .test
	go build -o ./bin/prof ./*.go

bin/prof-linux: GOOS=linux
bin/prof-linux: .test
	go build -o ./bin/prof-linux ./*.go

.docker: bin/prof-linux Dockerfile
	docker build -t dan353hehe/professor .
	@touch .docker
