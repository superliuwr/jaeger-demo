# stage 1) Build
FROM golang:1.17.6-alpine AS build-go

RUN apk add --no-cache \
            bash \
            curl \
            git \
            make && \
    rm -rf /var/cache/apk/*

WORKDIR /driver

# Add the rest of the source and build
COPY . /driver
RUN make build

# stage 2) Run
FROM golang:1.17.6-alpine

WORKDIR /app

COPY --from=build-go /driver/driver /app/

ENTRYPOINT ["./driver"]
