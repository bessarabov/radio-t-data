FROM golang:1.13.8-alpine3.11

RUN apk add --no-cache \
	git \
	sox \
	;

COPY main.go /go/src/

CMD go run src/main.go
