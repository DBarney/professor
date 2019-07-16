# super simplistic approach to detecting changes
SRC_FILES := $(shell find . -path ./.git -prune -o -name '*.go' -print)

# actual targets meant to be run
test: .test
image: .docker
build: bin/prof

.test: www/www.go internal/publisher/static.go $(SRC_FILES)
	go mod download
	go mod verify
	go vet ./...
	go test ./...
	@touch .test

www/www.go: www/index.html
	@echo rebuilding static html assets
	@go-bindata -fs -pkg www '-ignore=.*[.]go' -prefix www/ -o ./www/www.go ./www/...

internal/publisher/static.go: internal/publisher/error.md internal/publisher/failure.md internal/publisher/success.md
	@echo rebuilding markdown assets
	@go-bindata -fs -pkg publisher '-ignore=.*[.]go' -prefix internal/publisher/ -o ./internal/publisher/static.go ./internal/publisher/...
	

bin/prof: .test
	go build -o ./bin/prof ./*.go

bin/prof-linux: .test
	GOOS=linux go build -o ./bin/prof-linux ./*.go

.docker: bin/prof-linux Dockerfile
	docker build -t dan353hehe/professor .
	@touch .docker

