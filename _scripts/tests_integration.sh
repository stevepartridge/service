#!/usr/bin/env bash
set -e

BASE_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )/.." && pwd )"
CUR_DIR=$(pwd)

if [ -f "$BASE_DIR/_scripts/local.env" ]; then
  source "${BASE_DIR}/_scripts/local.env"
  # printenv
fi


cd $BASE_DIR

   
# printenv

go test ./tests -v $@


cd $CUR_DIR