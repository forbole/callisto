FROM golang:1.16 AS builder
WORKDIR /go/src/github.com/forbole/big-dipper
COPY . ./
RUN CGO_ENABLED=0 go build -o /build/bdjuno ./cmd/bdjuno

FROM alpine:latest
WORKDIR /root/
COPY --from=0 /build/bdjuno ./
RUN ./bdjuno init
ENTRYPOINT ["./bdjuno"]
CMD ["parse", "--home", "/root/.bdjuno"]
