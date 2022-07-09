FROM --platform=$BUILDPLATFORM golang:1.18.3 AS build
WORKDIR /rake/
COPY go.* .
COPY shared .
COPY server .
ARG PACKAGE="github.com/AppleGamer22/rake"
ARG VERSION="development"
ARG HASH="development"
ARG TARGETOS TARGETARCH
RUN GOOS=$TARGETOS GOARCH=$TARGETARCH go build -race -ldflags="-X '$PACKAGE/shared.Version=$VERSION' -X '$PACKAGE/shared.Hash=$HASH'" -o rake ./server

FROM --platform=$BUILDPLATFORM alpine:3.16 AS server
WORKDIR /rake/
COPY --from=build /rake/rake .
COPY templates .
ENV RAKE_STORAGE="/rake/storage"
ENV RAKE_DATABASE="rake"
EXPOSE 4100
CMD /rake/rake