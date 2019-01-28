FROM docker.io/liyinda/base_golang:1.9.7 AS build-env
MAINTAINER liyinda "liyinda@chinadaily.com.cn"

WORKDIR /go/src/

RUN cd /go/src \
  && apk update \
  && apk add git \
  && go get github.com/liyinda/AmrToMp3 \
  && go build -v  /go/src/github.com/liyinda/AmrToMp3/main.go

FROM docker.io/liyinda/base_alpine:3.6
COPY --from=build-env /go/src/main /usr/local/bin/amrtomp3
RUN chmod 775 /usr/local/bin/amrtomp3 \
  && apk update \
  && apk add tree tzdata \
  && cp -r -f /usr/share/zoneinfo/Hongkong /etc/localtime \
  && apk add yasm && apk add ffmpeg

CMD ["/usr/local/bin/amrtomp3"]
