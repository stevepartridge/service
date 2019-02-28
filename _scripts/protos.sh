#!/usr/bin/env bash

set -e

BASE_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )/.." && pwd )"
CUR_DIR=$(pwd)

PROJECT_PATH=github.com/stevepartridge/service/example

PROTOC_VERSION=3.6.1

if [[ "$1" == "setup" ]]; then

  if [ -d ./tmp ]; then 
    rm -Rf ./tmp
  fi

  mkdir -p ./tmp

  curl -f --ipv4 -Lo tmp/protoc-${PROTOC_VERSION}-osx-x86_64.zip \
    --connect-timeout 20 \
    --retry 6 \
    --retry-delay 10 \
    https://github.com/google/protobuf/releases/download/v${PROTOC_VERSION}/protoc-${PROTOC_VERSION}-osx-x86_64.zip

  cd tmp/
  unzip protoc-${PROTOC_VERSION}-osx-x86_64.zip
  chmod +x bin/protoc
  echo " Note: May need password to move protoc to /usr/local/bin/protoc"

  echo "Move protoc to /usr/local/bin ..."
  sudo cp bin/protoc /usr/local/bin/protoc
  echo "Copy include files to /usr/local/include ..."
  sudo cp -R include/google /usr/local/include/

  protoc --version
  
  # exit
  echo "Using go get to retreive grpc-ecosystem/grpc-gateway tools..."

  go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
  go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
  go get -u github.com/golang/protobuf/protoc-gen-go

  cd ../

  # this installs globally 
  if [ -d ./vendor/github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway ]; then
    # printf "vendor protoc-gen-grpc-gateway"
    go install ./vendor/github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
  elif [ -d $GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway ]; then
    # printf "gopath protoc-gen-grpc-gateway"
    go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
  else 
    echo "Unable to find github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway"; exit 1
  fi

  if [ -d ./vendor/github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger ]; then
    # printf "vendor protoc-gen-swagger"
    go install ./vendor/github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
  elif [ -d $GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger ]; then
    # printf "gopath protoc-gen-swagger"
    go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
  else
    echo "Unable to find github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger"; exit 1
  fi

  echo "Setup complete"
fi


cd $GOPATH/src/


printf "Go gRPC Files..."
protoc -I=./ \
  -I=$BASE_DIR/vendor/ \
  -I=$BASE_DIR/vendor/github.com/gogo/protobuf \
  -I=$BASE_DIR/vendor/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
  --go_out=plugins=grpc:. ./$PROJECT_PATH/protos/*.proto
  # --proto_path=. \

printf "Go gRPC Gateway Files..."
protoc -I=./ \
  -I=$BASE_DIR/vendor/ \
  -I=$BASE_DIR/vendor/github.com/gogo/protobuf \
  -I=$BASE_DIR/vendor/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
  --grpc-gateway_out=logtostderr=true:. ./$PROJECT_PATH/protos/*.proto

printf "Swagger/OpenAPI Files..."
protoc -I=./ \
  -I=$BASE_DIR/vendor/ \
  -I=$BASE_DIR/vendor/github.com/gogo/protobuf \
  -I=$BASE_DIR/vendor/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
  --proto_path=. --swagger_out=logtostderr=true:. ./$PROJECT_PATH/protos/*.proto


cd $CUR_DIR