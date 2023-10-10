# builder image
FROM golang:alpine3.16 as builder
RUN mkdir /build
ADD . /build/
WORKDIR /build
RUN go mod tidy
RUN go build -a -o fwallet .


# final image
FROM alpine:3.16
COPY --from=builder /build/fwallet .

EXPOSE 9090

# executable
ENTRYPOINT [ "./fwallet" ]