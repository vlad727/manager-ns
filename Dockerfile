FROM alpinelinux/golang AS builder
WORKDIR /app
COPY . /app
USER root
RUN env GOOS=linux GOARCH=amd64 && go build -o manager-ns /app/cmd/main.go

FROM alpine
WORKDIR /app
RUN apk update --no-check-certificate \
    && apk add --no-check-certificate curl net-tools
COPY --from=builder /app/manager-ns /app/manager-ns
RUN  chmod u+x manager-ns && mkdir /certs  /files
CMD ["./manager-ns"]