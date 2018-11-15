FROM golang:1.10.4 as build
WORKDIR /go/src/dotabot-ui

COPY telegram telegram
COPY handler handler
COPY vendor vendor
COPY state state
COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -a --installsuffix cgo -o dotabot-ui .

FROM alpine:3.8
WORKDIR /root

EXPOSE 8080

RUN apk --update --no-cache add ca-certificates

COPY --from=build /go/src/dotabot-ui/dotabot-ui .

COPY docker-entrypoint.sh .
RUN chmod +x docker-entrypoint.sh

CMD ["./docker-entrypoint.sh"]