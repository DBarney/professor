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

tree: a b c d e f g h i j k l m n o p q r s t u v w x y z

a c e g i k m o q s u w y:
	@echo $@
	@sleep 1

b d f h j l n p r t v x z:
	@echo $@