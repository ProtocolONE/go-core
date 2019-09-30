#!/usr/bin/env sh

if [ -n $1 ] && [ ${0:0:4} == "/bin" ]; then
  ROOT_DIR=$1/..
else
  ROOT_DIR="$( cd "$( dirname "$0" )" && pwd )/.."
fi

GO_PATH=${ROOT_DIR}/../../../..
GO_IMAGE=p1hub/go-core
GO_IMAGE_TAG=1.12
GO_PKG=github.com/ProtocolONE/go-core
GOOS="linux"
GOARCH="amd64"