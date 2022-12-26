# STAGE 1: BUILD
FROM golang:1.18.8-alpine3.16
ADD . /app
WORKDIR /app

# GET Signing CERTS
RUN apk update \
        && apk upgrade \
        && apk add --no-cache \
        ca-certificates \
        && update-ca-certificates 2>/dev/null || true

# DEP
RUN apk add --update --no-cache ca-certificates git

# VENDOR
RUN go mod download

# COMPILE
RUN mkdir -p ./bin
RUN GOGC=off CGO_ENABLED=0 GOOS=linux go build -gcflags=-trimpath=${GOPATH} -asmflags=-trimpath=${GOPATH} -a -o ./bin/api ./cmd/api/main.go
RUN GOGC=off CGO_ENABLED=0 GOOS=linux go build -gcflags=-trimpath=${GOPATH} -asmflags=-trimpath=${GOPATH} -a -o ./bin/migrate ./cmd/migrate/main.go

# STAGE 2: SCRATCH BINARY
FROM scratch
COPY --from=0 ./app/db /db
COPY --from=0 ./app/config /config
COPY --from=0 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=0 ./app/bin/api /bin/api
COPY --from=0 ./app/bin/migrate /bin/migrate

CMD ["/bin/api"]