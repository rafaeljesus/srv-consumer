FROM golang as builder
WORKDIR /go/src/github.com/rafaeljesus/srv-consumer
COPY . .
ENV CGO_ENABLED=0
ENV GOOS=linux
RUN make deps && make build

FROM scratch
WORKDIR /srv
COPY --from=builder /build/ .
CMD ["./srv-consumer"]
