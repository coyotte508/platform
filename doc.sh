#!/bin/sh -eu
# Generate OpenAPI documentation
GOPATH=${GOPATH:-~/go}
echo "Using GOPATH: ${GOPATH}"

if [ ! -x "$GOPATH/bin/swag" ]; then
  echo "Getting swag..."
  go get -u github.com/swaggo/swag/cmd/swag
fi

$GOPATH/bin/swag init -g ./services/data/data.go -o docs/api/v1/data