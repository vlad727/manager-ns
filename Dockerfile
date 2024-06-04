FROM alpinelinux/golang AS builder
WORKDIR /app
COPY . /app
USER root
RUN env GOOS=linux GOARCH=amd64 && go build -o webhook-app /app/cmd/main.go

FROM alpine
WORKDIR /app
COPY --from=builder /app/webhook-app /app/webhook-app
RUN  chmod u+x webhook-app && mkdir /certs  /files
CMD ["./webhook-app"]