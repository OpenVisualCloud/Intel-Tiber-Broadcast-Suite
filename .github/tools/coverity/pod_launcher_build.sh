#!/bin/bash

# script acording to readme
echo "**** BUILD pod Launcher ****"
cd ${1}/launcher/cmd/
go build main.go
cd ${1}