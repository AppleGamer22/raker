FROM --platform=$BUILDPLATFORM alpine:3.21.3 AS server
WORKDIR /raker
COPY raker .
# COPY ../../templates templates
# COPY ../../assets assets
RUN apk add ffmpeg
ENV STORAGE="/raker/storage"
ENV DATABASE="raker"
EXPOSE 4100
CMD ./raker