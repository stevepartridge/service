#!/usr/bin/env bash
# set -e

BASE_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )/.." && pwd )"

cd $BASE_DIR

if [ ! $(which go-bindata-assetfs) ]; then
	go get github.com/jteeuwen/go-bindata/...
  go get github.com/elazarl/go-bindata-assetfs/...
fi

go-bindata-assetfs -o "example/swagger.go" -pkg main -ignore=\\.sh -ignore=\\.go ./example/protos/...
go-bindata-assetfs -o "swagger/ui/swagger.go" -pkg swagger -ignore=\\.sh -ignore=\\.go ./swagger/ui/...
