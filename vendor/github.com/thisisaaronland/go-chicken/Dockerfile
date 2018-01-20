# https://blog.docker.com/2016/09/docker-golang/
# https://blog.golang.org/docker

# docker build -t rooster .
# docker run -it -p 1280:1280 rooster

# build phase - see also:
# https://medium.com/travis-on-docker/multi-stage-docker-builds-for-creating-tiny-go-images-e0e1867efe5a
# https://medium.com/travis-on-docker/triple-stage-docker-builds-with-go-and-angular-1b7d2006cb88

FROM golang:alpine AS build-env

RUN apk add --update alpine-sdk

ADD . /go-chicken

RUN cd /go-chicken; make bin

FROM alpine

COPY --from=build-env /go-chicken/bin/rooster /rooster

EXPOSE 8080

CMD /rooster -host 0.0.0.0
