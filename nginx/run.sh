#!/bin/bash

if [ ! "$(docker ps -aq -f name=nginx)" ]; then
	echo "gonna run"
	docker run --name nginx --hostname nginx --net=my-network -p 80:80 -d mynginx-image 
else
	if [ "$(docker ps -aq -f status=exited -f name=nginx)" ]; then
        docker start nginx
    fi
fi