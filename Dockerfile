FROM --platform=$BUILDPLATFORM golang:1.19.0-alpine AS build
WORKDIR /rake
COPY go.* .
COPY server server
COPY shared shared
ARG PACKAGE="github.com/AppleGamer22/rake"
ARG VERSION="development"
ARG HASH="development"
ARG TARGETOS TARGETARCH
ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV GOOS=$TARGETOS
ENV GOARCH=$TARGETARCH
RUN go build -ldflags="-X '$PACKAGE/shared.Version=$VERSION' -X '$PACKAGE/shared.Hash=$HASH'" -o rake ./server

FROM --platform=$BUILDPLATFORM alpine:3.16.1 AS server
WORKDIR /rake
COPY --from=build /rake/rake .
COPY templates templates
COPY assets assets
ENV STORAGE="/rake/storage"
ENV DATABASE="rake"
EXPOSE 4100
CMD ./rake