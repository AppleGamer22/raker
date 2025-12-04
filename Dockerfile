FROM alpine:3.23.0
WORKDIR /raker
COPY raker .
RUN apk add ffmpeg
ENV STORAGE="/raker/storage"
ENV DATABASE="raker"
EXPOSE 4100
CMD ./raker