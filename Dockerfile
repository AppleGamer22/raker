FROM --platform=$BUILDPLATFORM golang:1.22.6-alpine AS build
WORKDIR /raker
COPY go.* .
COPY server server
COPY shared shared
ARG PACKAGE="github.com/AppleGamer22/raker"
ARG VERSION="development"
ARG HASH="development"
ARG TARGETOS TARGETARCH
ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV GOOS=$TARGETOS
ENV GOARCH=$TARGETARCH
RUN go build -ldflags="-X '$PACKAGE/shared.Version=$VERSION' -X '$PACKAGE/shared.Hash=$HASH'" -o raker ./server

FROM --platform=$BUILDPLATFORM alpine:3.20.2 AS server
WORKDIR /raker
COPY --from=build /raker/raker .
COPY templates templates
COPY assets assets
RUN apk add ffmpeg
ENV STORAGE="/raker/storage"
ENV DATABASE="raker"
EXPOSE 4100
CMD ./raker