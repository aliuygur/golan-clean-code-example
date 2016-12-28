#!/usr/bin/env bash
./build.sh

GOPATH=`pwd`

pkill api
source .env
./bin/api &

inotifywait -m -r -e close_write src/app/ | while read line
do
	./build.sh && pkill api && ./bin/api &
done