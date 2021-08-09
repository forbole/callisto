FROM golang:1.16-alpine AS builder
RUN apk update && apk add --no-cache make git
WORKDIR /go/src/github.com/forbole/bdjuno
COPY . ./
RUN make build

FROM alpine:latest
RUN apk update && apk add --no-cache bash
WORKDIR /bdjuno
COPY --from=builder /go/src/github.com/forbole/bdjuno/build/bdjuno /usr/bin/bdjuno
COPY ./entrypoint.sh ./
ENTRYPOINT [ "./entrypoint.sh" ]
CMD [ "bdjuno" ]
