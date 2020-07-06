#!/bin/bash

if [ ! "$(docker ps -aq -f name=httpserver)" ]; then
	docker run --name  httpserver --hostname httpserver --net=my-network -p 8080:8080 -d httpserver-image 
else
	if [ "$(docker ps -aq -f status=exited -f name=httpserver)" ]; then
        docker start httpserver
    fi
fi