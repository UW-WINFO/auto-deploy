#!/usr/bin/env bash

# Use this script to (re)start the nginx web server that will serve
# the HTML content of the WINFO site.

if [[ "$(docker ps -aq --filter name=nginx)" ]]; then
            docker rm -f nginx
fi

docker run -d \
--restart unless-stopped \
-p 80:80 \
-p 443:443 \
-v /site/Website/:/usr/share/nginx/html:ro \
-v /etc/letsencrypt/live/winfo.ischool.uw.edu/fullchain.pem:/etc/pki/tls/certs/winfo.ischool.uw.edu-cert.pem \
-v /etc/letsencrypt/live/winfo.ischool.uw.edu/privkey.pem:/etc/pki/tls/private/winfo.ischool.uw.edu-key.pem \
--name nginx \
winfo/website
