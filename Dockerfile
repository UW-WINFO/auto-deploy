FROM alpine
RUN apk add --no-cache ca-certificates
RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh

ADD auto-deploy auto-deploy

ENTRYPOINT [ "/auto-deploy" ]