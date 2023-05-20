FROM botwayorg/pb-core:latest AS download
FROM golang:1.20.4 as build

WORKDIR /go/src/app
COPY . .

RUN go mod download

RUN CGO_ENABLED=0 go build -o /go/bin/app

FROM debian:bullseye-slim


ARG AUTH_TOKEN
ARG REPO="__pocketbasedata"

ENV AUTH_TOKEN=$AUTH_TOKEN
ENV REPO=$REPO

ENV DIR "/root/pocketbase"
ENV CMD "/usr/local/bin/pocketbase serve --http=0.0.0.0:8090 --dir=/root/pocketbase"

RUN apt-get update && apt install -y git git-lfs

COPY --from=download /pocketbase /usr/local/bin/pocketbase
COPY --from=build /go/bin/app /usr/local/bin/app

EXPOSE 8090

RUN mkdir -p /root/pocketbase

ENTRYPOINT [ "app" ]