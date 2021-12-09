#!/bin/bash

docker run -v ${PWD}/nginx.conf:/etc/nginx/nginx.conf -p 8007:80 nginx:1.17.1