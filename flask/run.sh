#!/bin/bash

if [ ! "$(docker ps -aq -f name=flask)" ]; then
	docker run --name flask --hostname flask --net=my-network -p 5000:5000 -d myflask-image
else
	if [ "$(docker ps -aq -f status=exited -f name=flask)" ]; then
        docker start flask
    fi
fi