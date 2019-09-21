set -e

ROOT_DIR=$(dirname $0)

pushd $ROOT_DIR

rm -rf ./build
mkdir -p ./build
go get ./pkg/*
go test ./pkg/*
GOOS=linux go build -o ./build/one-time-secret-http-linux pkg/http/http.go
GOOS=darwin go build -o ./build/one-time-secret-http-darwin pkg/http/http.go

GOOS=linux GOARCH=amd64 go build -o ./build/one-time-secret-index pkg/lambda-index/main.go
GOOS=linux GOARCH=amd64 go build -o ./build/one-time-secret-create pkg/lambda-create/main.go

zip -j ./build/one-time-secret-index.zip ./build/one-time-secret-index
zip -j ./build/one-time-secret-create.zip ./build/one-time-secret-create

if [[ "$1" == "--deploy" ]]; then
  pushd terraform/root
  terraform apply
  popd
fi

popd
