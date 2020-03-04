FROM golang:1.13 as build

RUN go env -w GOPROXY=https://goproxy.cn,direct

WORKDIR /go/cache
ADD . .
RUN go mod download

RUN GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -installsuffix cgo -o=pegasus_engine -tags=jsoniter cmd/pegasus_engine.go

FROM alpine:3.9 as prod
COPY --from=build /go/cache/pegasus_engine /
CMD ["/pegasus_engine"]
