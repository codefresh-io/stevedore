FROM golang:1.19-alpine3.16 as builder

# Add basic tools
RUN apk add --no-cache --update curl bash make git

RUN mkdir -p /go/src/github.com/codefresh-io/stevedore
WORKDIR /go/src/github.com/codefresh-io/stevedore

ENV GOPATH /go

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

# Build binary
RUN "./scripts/BUILD.sh"


FROM alpine:3.18.9

RUN apk add --no-cache ca-certificates

COPY --from=builder /go/src/github.com/codefresh-io/stevedore/dist/bin/stevedore /usr/bin/stevedore


ENV PATH $PATH:/usr/bin/stevedore
ENTRYPOINT ["stevedore"]

CMD ["--help"]