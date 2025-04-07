FROM golang:1.24.0-alpine AS build

RUN mkdir /build
WORKDIR /build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o status-webhooks


FROM alpine
ENV TZ=Asia/Taipei
ENV LC_ALL=C.UTF-8
ENV LANG=C.UTF-8
RUN apk update --no-cache && apk --no-cache add gcc g++ make ca-certificates \
  && apk add git curl  net-tools \
  && apk add -U tzdata \
  && cp /usr/share/zoneinfo/$TZ /etc/localtime \
  && echo $TZ > /etc/timezone 
# && apk del tzdata
RUN addgroup -S appUser \
  && adduser -S -D appUser appUser
USER appUser
COPY --from=build --chown=appUser:appUser /build/hls-key-server-go /app/
COPY --from=build --chown=appUser:appUser /build/config /app/config/
WORKDIR /app
ENTRYPOINT ["./hls-key-server-go"]
