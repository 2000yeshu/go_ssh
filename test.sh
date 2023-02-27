#! /bin/bash

sudo docker build -t test_ssh_server -f ./Dockerfile_ssh .
CONTAINER_ID=$(sudo docker run -d -P test_ssh_server)
sudo docker inspect --format='{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' $CONTAINER_ID

go run main.go --containerid $CONTAINER_ID