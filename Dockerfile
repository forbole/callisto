FROM golang:1.16-alpine AS builder
RUN set -eux; apk add --no-cache ca-certificates build-base;
RUN apk add git

WORKDIR /go/src/github.com/forbole/bdjuno
COPY . ./

RUN make docker-build

FROM alpine:latest
WORKDIR /bdjuno
COPY --from=builder /go/src/github.com/forbole/bdjuno/build/bdjuno /usr/bin/bdjuno
COPY --from=builder /go/src/github.com/forbole/bdjuno/volume /bdjuno
CMD [ "bdjuno" ]
