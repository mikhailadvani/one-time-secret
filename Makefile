all: build docker

install:
	dep ensure

test:
	go test ./pkg/...

build: build_local build_lambda

build_local: install build_local_linux build_local_darwin

build_local_linux:
	mkdir -p ./build
	GOOS=linux go build -o ./build/one-time-secret-http-linux pkg/http/http.go

build_local_darwin:
	mkdir -p ./build
	GOOS=darwin go build -o ./build/one-time-secret-http-darwin pkg/http/http.go

build_lambda: install
	mkdir -p ./build/lambda
	GOOS=linux GOARCH=amd64 go build -o ./build/lambda/one-time-secret-index pkg/lambda-index/main.go
	GOOS=linux GOARCH=amd64 go build -o ./build/lambda/one-time-secret-create pkg/lambda-create/main.go
	GOOS=linux GOARCH=amd64 go build -o ./build/lambda/one-time-secret-get pkg/lambda-get/main.go
	zip -j ./build/one-time-secret-index.zip ./build/lambda/one-time-secret-index
	zip -j ./build/one-time-secret-create.zip ./build/lambda/one-time-secret-create
	zip -j ./build/one-time-secret-get.zip ./build/lambda/one-time-secret-get
	rm -rf ./build/lambda

clean:
	rm -rf ./build

run_local:
	go run ./pkg/http/http.go

run_lambda: build_lambda
	terraform init terraform/root
	terraform apply -auto-approve terraform/root

destroy_lambda:
	terraform init terraform/root
	terraform destroy -auto-approve terraform/root

docker: install
	docker build -t one-time-secret .

fmt:
	go fmt ./pkg/...
