FROM golang:1.18-alpine AS builder

ARG GIT_TOKEN

RUN go env -w GOPRIVATE=github.com/NibiruChain
RUN git config --global url."https://git:${GIT_TOKEN}@github.com".insteadOf "https://github.com"

RUN apk update && apk add --no-cache make git
WORKDIR /go/src/github.com/forbole/bdjuno
COPY . ./
RUN go mod download
RUN make build

FROM gcr.io/distroless/base
WORKDIR /bdjuno
COPY --from=builder /go/src/github.com/forbole/bdjuno/build/bdjuno /usr/bin/bdjuno
ENTRYPOINT [ "bdjuno" ]
