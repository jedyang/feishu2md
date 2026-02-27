ARG GO_VERSION=1.24
FROM docker.cnb.cool/meyley/docker/golang:${GO_VERSION}-alpine AS builder

WORKDIR /feishu2md

COPY go.mod go.sum ./
RUN export GOPROXY=https://goproxy.cn,direct && go mod download

COPY core  ./core
COPY web ./web
COPY utils ./utils
RUN go build -o ./feishu2md4web ./web/*.go

FROM docker.cnb.cool/meyley/docker/alpine:latest
RUN apk update && apk add --no-cache ca-certificates

ENV GIN_MODE=release

COPY --from=builder /feishu2md/feishu2md4web ./

EXPOSE 8080

ENTRYPOINT ["./feishu2md4web"]
