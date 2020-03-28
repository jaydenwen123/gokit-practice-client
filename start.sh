#!/bin/bash
echo "begin to complie the program..."
client="myClient"
go build -o $client
echo "finish the complie...."
echo "begin to start exec the client...."
./$client
