FROM amd64/golang:1.17-buster AS builder
USER root

RUN apt-get update && apt-get install git
# RUN curl -L https://github.com/hasura/graphql-engine/raw/stable/cli/get.sh | bash
WORKDIR /go/src/github.com/forbole/bdjuno
COPY . ./
RUN go mod tidy -compat=1.17
RUN make build
RUN FOLDER=$(ls /go/pkg/mod/github.com/\!cosm\!wasm/ | grep wasmvm@v) && ln -s /go/pkg/mod/github.com/\!cosm\!wasm/${FOLDER} /go/pkg/mod/github.com/\!cosm\!wasm/wasmvm


FROM amd64/golang:1.17-buster
USER root

WORKDIR /bdjuno
#TODO: fix line 18
#COPY --from=builder /go/pkg/mod/github.com/!cosm!wasm/wasmvm/api/libwasmvm.so /usr/lib
COPY --from=builder /go/src/github.com/forbole/bdjuno/build/bdjuno /usr/bin/bdjuno
# COPY --from=builder /usr/local/bin/hasura /usr/local/bin/hasura
COPY --from=builder /go/src/github.com/forbole/bdjuno/hasura /hasura
COPY bdjuno/ /usr/local/bdjuno/bdjuno/

ARG HASURA_ACTIONS_PORT
ARG HASURA_ACTIONS_GRPC
ARG HASURA_ACTIONS_RPC

ENV HASURA_ACTIONS_PORT ${HASURA_ACTIONS_PORT}
ENV HASURA_ACTIONS_GRPC ${HASURA_ACTIONS_GRPC}
ENV HASURA_ACTIONS_RPC ${HASURA_ACTIONS_RPC}

CMD ["/bin/bash", "-c", "bdjuno hasura-actions --port ${HASURA_ACTIONS_PORT} --grpc ${HASURA_ACTIONS_GRPC} --rpc ${HASURA_ACTIONS_RPC} --home /usr/local/bdjuno/bdjuno/"]
