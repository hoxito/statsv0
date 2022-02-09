# Docker para desarrollo
FROM golang:1.14.3-buster

WORKDIR /go/src/statsv0

ENV REDIS_URL host.docker.internal:6379
ENV RABBIT_URL amqp://host.docker.internal
ENV AUTH_SERVICE_URL http://host.docker.internal:3000

# Puerto de stats Service y debug
EXPOSE 3010

CMD ["go" , "run" , "/go/src/statsv0"]