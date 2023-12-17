FROM golang:1.21.5 as builder

ENV CGO_ENABLED 0

COPY . /client

WORKDIR /client/cmd/client
RUN go build

FROM alpine:3.18

RUN addgroup -g 1000 -S client && \
    adduser -u 1000 -h /client -G client -S client

COPY --from=builder --chown=client:client /client/cmd/client/client /client/client

WORKDIR /client

USER client

CMD ["./client", "--host=server"]


