#!/usr/bin/env bash

GOOS=windows GOARCH=amd64 go build -trimpath -ldflags "-s -w -h -H windowsgui -extldflags=-static" -o bin/multibuild-win-amd64.exe .
GOOS=linux GOARCH=amd64 go build -trimpath -ldflags "-s -w -h -extldflags=-static" -o bin/multibuild-linux-amd64 .
