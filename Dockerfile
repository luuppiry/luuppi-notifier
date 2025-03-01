FROM	golang:1.24 AS builder

WORKDIR	/usr/src/app
COPY 	. .
RUN 	go mod download
RUN 	make build
RUN 	make test

FROM 	gcr.io/distroless/static-debian12:latest
USER 	1000:1000
COPY 	--from=builder --chown=1000:1000 --chmod=700 /usr/src/app/build/notifier /notifier
COPY 	--chown=1000:1000 --chmod=600 config.json /config.json 
ENTRYPOINT ["/notifier"]
