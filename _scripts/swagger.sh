#!/usr/bin/env bash
# set -e

BASE_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )/.." && pwd )"

cd $BASE_DIR

# Setup
if [ "$1" == "setup" ]; then
	git clone git@github.com:swagger-api/swagger-ui.git ./tmp/

	cd ./tmp 
	git checkout v3.17.2

	if [ -d ./swagger/ui ]; then
		rm -Rf ./swagger/ui
	fi 

	cd ../

	if [ ! -d ./swagger/ui ]; then
		mkdir -p ./swagger/ui
	fi 
	mv ./tmp/dist ./swagger/ui

	rm -Rf ./tmp

fi

if [ ! $(which go-bindata-assetfs) ]; then
	go get github.com/jteeuwen/go-bindata/...
  go get github.com/elazarl/go-bindata-assetfs/...
fi

go-bindata-assetfs -o "swagger/ui.go" -pkg swagger -ignore=\\.sh -ignore=\\.go ./swagger/ui/...
