FROM golang:1.21.5 as builder

ENV CGO_ENABLED 0

COPY . /server

WORKDIR /server/cmd/server
RUN go build

FROM alpine:3.18

RUN addgroup -g 1000 -S server && \
    adduser -u 1000 -h /server -G server -S server

COPY --from=builder --chown=server:server /server/cmd/server/server /server/server

WORKDIR /server

USER server

EXPOSE 3333

CMD ["./server", "--host=0.0.0.0"]


