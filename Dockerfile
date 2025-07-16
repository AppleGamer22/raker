FROM alpine:3.22.1
WORKDIR /raker
COPY raker .
RUN apk add ffmpeg
ENV STORAGE="/raker/storage"
ENV DATABASE="raker"
EXPOSE 4100
CMD ./raker