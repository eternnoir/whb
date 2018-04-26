#!/bin/bash
go fmt `go list ./... | grep -v vendor`
go get golang.org/x/tools/cmd/goimports
goimports -w -d $(find . -type f -name '*.go' -not -path "./vendor/*")
