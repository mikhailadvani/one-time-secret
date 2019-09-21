
install:
	dep ensure

test:
	go test ./pkg/...

build:
	mkdir -p ./build
	mkdir -p ./build/lambda
	GOOS=linux go build -o ./build/one-time-secret-http-linux pkg/http/http.go
	GOOS=darwin go build -o ./build/one-time-secret-http-darwin pkg/http/http.go
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

run_lambda: build
	cd terraform/root
	terraform init
	terraform apply -auto-approve

fmt:
	go fmt ./pkg/... ./test/...
