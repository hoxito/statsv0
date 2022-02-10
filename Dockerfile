# Docker para desarrollo
FROM golang:1.14.3-buster

WORKDIR /go/src/statsv0

ENV REDIS_URL host.docker.internal:6379
ENV RABBIT_URL amqp://host.docker.internal
ENV AUTH http://host.docker.internal:3000
ENV CATALOG http://host.docker.internal:3002
ENV ORDERS http://host.docker.internal:3004
ENV MONGO_URL mongodb+srv://jose:statsgo@cluster0.j1j5b.mongodb.net/myFirstDatabase?retryWrites=true&w=majority
# Puerto de stats Service y debug
EXPOSE 3010

CMD ["go" , "run" , "/go/src/statsv0"]