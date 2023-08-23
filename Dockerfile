FROM golang:1.16-alpine AS builder
RUN set -eux; apk add --no-cache ca-certificates build-base;
RUN apk add git

WORKDIR /go/src/github.com/forbole/bdjuno
COPY . ./

# See https://github.com/CosmWasm/wasmvm/releases
ADD https://github.com/CosmWasm/wasmvm/releases/download/v0.16.1/libwasmvm_muslc.a /lib/libwasmvm_muslc.a
RUN sha256sum /lib/libwasmvm_muslc.a | grep 0e62296b9f24cf3a05f8513f99cee536c7087079855ea6ffb4f89b35eccdaa66

RUN make docker-build

FROM alpine:latest
WORKDIR /bdjuno
COPY --from=builder /go/src/github.com/forbole/bdjuno/build/bdjuno /usr/bin/bdjuno
COPY --from=builder /go/src/github.com/forbole/bdjuno/volume /bdjuno
CMD [ "bdjuno" ]
