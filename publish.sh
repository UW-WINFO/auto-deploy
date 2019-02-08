#!/usr/bin/env bash
echo 'Building auto-deploy linux binary...';
GOOS=linux go build

# Going forward, if someone from WINFO wants to make modifications
# to this we should create a WINFO docker hub organization and push
# this container to that so WINFO can have push access.
echo 'Buidling auto-deploy docker container...';
docker build -t brendankellogg/winfo-deploy-server .
echo 'Pushing auto-deploy docker container...';
docker push brendankellogg/winfo-deploy-server

echo 'Cleaning up...'
go clean
echo 'Finished'
