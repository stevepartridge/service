#!/usr/bin/env bash
# set -e

BASE_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )/.." && pwd )"

VERSION=3.46.0

cd $BASE_DIR

tmpDir=$BASE_DIR/tmp
[ ! -d $tmpDir ] || rm -Rf $tmpDir

mkdir -p $tmpDir

cd $tmpDir

curl -L https://github.com/swagger-api/swagger-ui/archive/refs/tags/v${VERSION}.tar.gz | tar zx

mv swagger-ui-${VERSION}/dist/* $BASE_DIR/static/.

cd $BASE_DIR

rm -Rf $tmpDir

ls -al