#!/usr/bin/env bash

# Use this script to (re)start the deploy hook server that will
# listen for incoming hooks from Github and update the local content
# as necessary.

docker pull brendankellogg/winfo-deploy-server
if [[ "$(docker ps -aq --filter name=deploy-server)" ]]; then
            docker rm -f deploy-server
fi
docker run -d -p 4000:4000 \
    --restart unless-stopped \
        -v /etc/letsencrypt/live/winfo.ischool.uw.edu/fullchain.pem:/etc/pki/tls/certs/winfo.ischool.uw.edu-cert.pem:ro \
        -v /etc/letsencrypt/live/winfo.ischool.uw.edu/privkey.pem:/etc/pki/tls/private/winfo.ischool.uw.edu-key.pem:ro \
	-v /site/Website/:/site/Website/ \
	-e ADDR=:4000 \
        -e TLSKEY=/etc/pki/tls/private/winfo.ischool.uw.edu-key.pem \
        -e TLSCERT=/etc/pki/tls/certs/winfo.ischool.uw.edu-cert.pem \
	-e AUTO_UPDATE_CONTENT_DIR=/site/Website \
	-e AUTO_UPDATE_GIT_REPO=https://github.com/UW-WINFO/Website.git \
	--name deploy-server \
	brendankellogg/winfo-deploy-server
