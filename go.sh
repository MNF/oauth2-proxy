#!/bin/bash

# https://medium.com/iron-io-blog/the-easiest-way-to-develop-with-go-introducing-a-docker-based-go-tool-c456238507d6#.atcfnlf60
if [ "$1" == "run" ]; then  
  shift
  docker run -it --rm -v $PWD:/app -w /app -p 8080:8080 iron/base ./app $@
elif [ "$1" == "bash" ]; then
  shift
  docker run -it --rm -v $PWD:/app -w /app -p 8080:8080 iron/base /bin/ash $@
else # vendor, build 
  docker run -it --rm -v $PWD:/app -w /app treeder/go $@
fi
