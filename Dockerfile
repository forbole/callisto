FROM golang:1.17-buster AS builder
RUN apt-get update && apt-get install git
WORKDIR /go/src/github.com/forbole/bdjuno
RUN git clone https://github.com/CudoVentures/cosmos-gravity-bridge.git ../CudosGravityBridge
COPY . ./
RUN make build


FROM golang:1.17-buster
WORKDIR /bdjuno
COPY --from=builder /go/pkg/mod/github.com/!cosm!wasm/wasmvm@v0.16.0/api/libwasmvm.so /usr/lib
COPY --from=builder /go/src/github.com/forbole/bdjuno/build/bdjuno /usr/bin/bdjuno
# COPY config.yaml /usr/local/bdjuno/.bdjuno/
CMD [ "bdjuno parse --home /usr/local/bdjuno/.bdjuno/"]
# CMD ["/bin/bash", "-c", "sleep infinity"]

