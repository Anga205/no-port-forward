#!/bin/bash

cd tests/

go run service.go &

GO_PID=$!

sleep 1

httpResponse=$(curl --location http://127.0.0.1:8080)

echo $httpResponse

kill $GO_PID
