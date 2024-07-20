#!/bin/bash

cd tests/

gcc service.c -o service

./service &

HTTP_PID=$!

PORT=6000

nc -z 127.0.0.1 $PORT

if [ $? -eq 0 ]; then
    echo "Port $PORT is open"
    httpResponse=$(curl --location http://127.0.0.1:$PORT)
    echo $httpResponse
else
    echo "Port $PORT is closed"
fi

kill $HTTP_PID
