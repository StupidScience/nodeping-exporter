FROM alpine:3.9

EXPOSE 9503

RUN apk add --no-cache ca-certificates

COPY nodeping-exporter ./

ENTRYPOINT ["./nodeping-exporter"]
