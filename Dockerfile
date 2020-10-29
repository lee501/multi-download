FROM golang:latest as build
ENV GOPROXY https://goproxy.io, direct
ENV GO111MODULE on
WORKDIR /go/cache
ADD go.mod .
ADD go.sum .
RUN go mod download

WORKDIR /go/release
ADD . .
RUN GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -installsuffix cgo -o server .

#FROM scratch as prod
#COPY --from=build /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
#COPY --from=build /go/release/server /

FROM debian:9
RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai  /etc/localtime

RUN apt-get update && apt-get install -y --no-install-recommends \
  ca-certificates \
  && update-ca-certificates \
  && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY --from=build /go/release/server .
