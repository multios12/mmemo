#!/bin/sh

export HTML="index.html"
yarn --cwd ./front build
cp ./front/dist/* ./cmd/mmemo/static/ -R

cd cmd/mmemo
export GOOS=linux
export GOARCH=amd64
go build -ldflags="-s -w" -trimpath -o ../../dist/
cd ..