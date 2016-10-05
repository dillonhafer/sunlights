#!/bin/bash
set -e
go test
GOOS=linux GOARCH=arm GOARM=6 go build github.com/dillonhafer/sunlights
