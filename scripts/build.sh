#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

if [ -z "${BIN}" ]; then
    echo "BIN must be set"
    exit 1
fi
if [ -z "${PKG}" ]; then
    echo "PKG must be set"
    exit 1
fi
if [ -z "${VERSION}" ]; then
    echo "VERSION must be set"
    exit 1
fi

export CGO_ENABLED=0

echo "Go installing app with PKG: ${PKG} VERSION: ${VERSION}"
go build -v                                             \
    -installsuffix "static"                             \
    -ldflags "-X ${PKG}/pkg/version.VERSION=${VERSION}" \
	-o bin/${BIN}                                       \
    ${PKG}/cmd/${BIN}/...
echo "Successfully installed, exiting build"
