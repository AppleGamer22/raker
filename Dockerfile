FROM --platform=$BUILDPLATFORM golang:1.18.3-alpine AS builder
WORKDIR /rake
COPY . .
ARG TARGETOS TARGETARCH
RUN GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o rake ./server

FROM --platform=$BUILDPLATFORM alpine:3.16
WORKDIR /rake
COPY --from=builder /rake/rake .
ENV RAKE_STORAGE "/rake/storage"
ENV RAKE_DATABASE "rake"
ENV RAKE_PORT 4100
EXPOSE 4100
CMD /rake/rake