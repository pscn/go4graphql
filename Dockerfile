FROM golang:1.12

# Prepare env
ENV GOPATH /go
ENV PATH ${GOPATH}/bin:$PATH

# Install dep
RUN go get -u github.com/golang/dep/cmd/dep

RUN mkdir -p /go/src/github.com/pscn/go4graphql
WORKDIR /go/src/github.com/pscn/go4graphql
COPY . .

RUN dep ensure
RUN go build

CMD ["./go4graphql"]
