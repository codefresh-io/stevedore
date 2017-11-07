FROM golang:latest as builder
RUN mkdir -p /go/src/github.com/codefresh-io/stevedore
WORKDIR /go/src/github.com/codefresh-io/stevedore
COPY . .
RUN "./scripts/BUILD.sh"


FROM alpine:3.6

RUN apk add --no-cache ca-certificates

COPY --from=builder /go/src/github.com/codefresh-io/stevedore/dist/bin/stevedore /usr/bin/stevedore
ENV PATH $PATH:/usr/bin/stevedore
ENTRYPOINT ["stevedore"]

CMD ["--help"]