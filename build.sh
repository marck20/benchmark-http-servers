#!/bin/bash

docker build -t myflask-image flask/

docker build -t mynginx-image nginx/

docker build -t httpserver-image httpserver/

docker build -t ab-image ab/

docker build -t goab_channel-image goab/channel/

docker build -t goab_mutex-image goab/mutex/

docker build -t goab_channel_wp-image goab/channel_wp/

docker build -t goab_mutex_wp-image goab/mutex_wp/



