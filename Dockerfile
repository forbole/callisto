FROM golang:1.16-alpine AS builder
RUN apk update && apk add --no-cache make git
WORKDIR /go/src/github.com/forbole/bdjuno
COPY . ./
RUN go mod download
RUN make build

FROM gcr.io/distroless/base
WORKDIR /bdjuno
COPY --from=builder /go/src/github.com/forbole/bdjuno/build/bdjuno /usr/bin/bdjuno
ENTRYPOINT [ "bdjuno" ]
