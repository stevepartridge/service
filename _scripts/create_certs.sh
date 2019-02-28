#!/usr/bin/env bash
set -o errexit
# set -o nounset
set -o pipefail

BASE_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )/.." && pwd )"
CUR_DIR=$(pwd)

if [ ! -d $BASE_DIR/certificates ]; then
	mkdir -p $BASE_DIR/certificates
fi
cd $BASE_DIR/certificates

if [[ -d out && "$1" == "replace" ]]; then
  rm -Rf out
fi

ORG="Local Development Corp. LLC. Inc."
ORG_UNIT="Local Development"
COMMON_NAME="host.local"
DOMAIN="*.host.local,localhost,127.0.0.1"
EXPIRES="2 years"
COUNTRY="US"

if [[ ! $(which certstrap) ]]; then 
	echo ""
	echo "Install certstrap first, then try again"
	echo ""
	echo "  go get -u github.com/square/certstrap"
	echo ""
	exit 1
fi

echo "=== init cert authority ==="
certstrap init --common-name "${ORG} CA" \
  --expires "${EXPIRES}" \
  --organization "${ORG}" \
  --organizational-unit "${ORG_UNIT}" \
  --country "${COUNTRY}" \
  --province "California" \
  --locality "San Diego" \
  --passphrase ""
  # --depot-path "$BASE_DIR/certificates"

echo "=== request certs ==="
certstrap request-cert --common-name "${COMMON_NAME}" \
  --ip 127.0.0.1 \
  --domain ${DOMAIN} \
  --organization "${ORG}" \
  --organizational-unit "${ORG_UNIT}" \
  --country "${COUNTRY}" \
  --province "California" \
  --locality "San Diego" \
  --passphrase ""

echo "=== sign certs ==="
certstrap sign "${COMMON_NAME}" \
  --CA "${ORG} CA" \
  --expires "${EXPIRES}" \
  --passphrase ""
  # --intermediate 
echo "=== done ==="

# mv $BASE_DIR/out $BASE_DIR/certificates

CERT=$(cat $BASE_DIR/certificates/out/${COMMON_NAME}.crt)
KEY=$(cat $BASE_DIR/certificates/out/${COMMON_NAME}.key)
ROOTCA=$(cat $BASE_DIR/certificates/out/${ORG// /_}_CA.crt)

certsgo=$(cat <<EOF
package insecure

const (
  Cert = \`${CERT}\`
  Key = \`${KEY}\`
  RootCA = \`${ROOTCA}\`
)

EOF
)

echo "$certsgo" > $BASE_DIR/insecure/certs.go

cd $CUR_DIR