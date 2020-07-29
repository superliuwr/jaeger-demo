# stage 1) Build
FROM golang:1.13.4-alpine AS build-go

RUN apk add --no-cache \
            bash \
            curl \
            git \
            make && \
    rm -rf /var/cache/apk/*

WORKDIR /frontend

# Add the rest of the source and build
COPY . /frontend
RUN make build

# stage 2) Run
FROM golang:1.13.4-alpine

WORKDIR /app

COPY --from=build-go /frontend/frontend /app/

ENTRYPOINT ["./frontend"]
