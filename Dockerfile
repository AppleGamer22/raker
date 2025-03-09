FROM alpine:3.21.3
WORKDIR /raker
COPY raker .
# COPY ../../templates templates
# COPY ../../assets assets
RUN apk add ffmpeg libc6-compat
ENV STORAGE="/raker/storage"
ENV DATABASE="raker"
EXPOSE 4100
CMD ./raker