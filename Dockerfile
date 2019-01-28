FROM golang:1.11.4 as feugo-server
WORKDIR /go/src/github.com/sridharavinash/feugo
COPY . .
RUN go build -o bin/server

FROM debian:stretch
EXPOSE 8081
COPY --from=feugo-server /go/src/github.com/sridharavinash/feugo/bin/server /
COPY --from=feugo-server /go/src/github.com/sridharavinash/feugo/assets /assets
COPY --from=feugo-server /go/src/github.com/sridharavinash/feugo/public /public

ENTRYPOINT ["/server"]
