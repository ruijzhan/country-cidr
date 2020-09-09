#!/usr/bin/env bash

#go get -u github.com/go-bindata/go-bindata/...
GEN_DATA_FILE=1 go test -run Test_ParseAPNIC -count=1
go-bindata -pkg country_cidr -o data.go apnic.json
rm -f apnic.json
