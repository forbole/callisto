FROM golang:1.17-buster AS builder
RUN apt-get update && apt-get install git
WORKDIR /go/src/github.com/forbole/bdjuno
RUN git clone https://github.com/CudoVentures/cosmos-gravity-bridge.git ../CudosGravityBridge
COPY . ./
RUN go mod tidy
RUN make build
RUN FOLDER=$(ls /go/pkg/mod/github.com/\!cosm\!wasm/ | grep wasmvm@v) && ln -s /go/pkg/mod/github.com/\!cosm\!wasm/${FOLDER} /go/pkg/mod/github.com/\!cosm\!wasm/wasmvm


FROM golang:1.17-buster
WORKDIR /bdjuno
COPY --from=builder /go/pkg/mod/github.com/!cosm!wasm/wasmvm/api/libwasmvm.so /usr/lib
COPY --from=builder /go/src/github.com/forbole/bdjuno/build/bdjuno /usr/bin/bdjuno
COPY bdjuno/ /usr/local/bdjuno/bdjuno/
CMD ["/bin/bash", "-c", "bdjuno parse --home /usr/local/bdjuno/bdjuno/"]
# CMD ["/bin/bash", "-c", "sleep infinity"] 

