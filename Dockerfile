FROM golang:alpine AS builder
LABEL maintainer="infra-team@cardinals"

RUN apk add --update git bash openssh make gcc musl-dev

WORKDIR /go/src/Cardinal-Cryptography/github-actions-runners-exporter
COPY . .
RUN go build

FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /bin
COPY --from=builder /go/src/Cardinal-Cryptography/github-actions-runners-exporter/github-actions-runners-exporter github-actions-runners-exporter
RUN chmod +x /bin/github-actions-runners-exporter
RUN /bin/github-actions-runners-exporter
ENTRYPOINT ["/bin/github-actions-runners-exporter"]
