#!/bin/bash

GOOS=windows GOARCH=amd64 go build  -o SASC-Win64.exe SASC.go
GOOS=linux   GOARCH=amd64 go build  -o SASC-Linux64 SASC.go
GOOS=darwin  GOARCH=amd64 go build  -o SASC-MacOS SASC.go

