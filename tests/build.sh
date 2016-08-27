#!/usr/bin/env bash
set -euo pipefail
IFS=$'\n\t'

cd $GOPATH/src/github.com/variadico/noti/cmd/noti || exit 1

echo "# build test"; echo

echo "## darwin"
GOOS=darwin go build && rm noti
echo "## linux"
GOOS=linux go build && rm noti
echo "## freebsd"
GOOS=freebsd go build && rm noti
echo "## windows"
GOOS=windows go build && rm noti.exe
