#!/bin/sh

echo "Building stevedore"
CGO_ENABLED=0 go build -v -o "./dist/bin/stevedore" *.go