FROM alpine:3.23.4

# target platform is provided by GoReleaser's docker build context
ARG TARGETPLATFORM

WORKDIR /raker

# copy the platform-specific binary produced by GoReleaser
COPY $TARGETPLATFORM/raker .

# copy bundled frontend assets
COPY vdist vdist

# install runtime dependencies
RUN apk add --no-cache ffmpeg

ENV STORAGE="/raker/storage"
ENV DATABASE="raker"

EXPOSE 4100
CMD ["./raker"]