FROM amd64/golang:1.17-buster AS builder
USER root

RUN apt-get update && apt-get install git
WORKDIR /go/src/github.com/forbole/bdjuno
COPY . ./
RUN go mod tidy
RUN make build
RUN FOLDER=$(ls /go/pkg/mod/github.com/\!cosm\!wasm/ | grep wasmvm@v) && ln -s /go/pkg/mod/github.com/\!cosm\!wasm/${FOLDER} /go/pkg/mod/github.com/\!cosm\!wasm/wasmvm


FROM amd64/golang:1.17-buster
USER root

WORKDIR /bdjuno
COPY --from=builder /go/pkg/mod/github.com/!cosm!wasm/wasmvm/api/libwasmvm.so /usr/lib
COPY --from=builder /go/src/github.com/forbole/bdjuno/build/bdjuno /usr/bin/bdjuno
COPY --from=builder /go/src/github.com/forbole/bdjuno/hasura /hasura
COPY bdjuno/ /usr/local/bdjuno/bdjuno/

CMD ["/bin/bash", "-c", "bdjuno parse --home /usr/local/bdjuno/bdjuno/"]
