FROM golang:1.18-alpine as build
WORKDIR /app
ADD go.mod go.sum ./
ADD . .
RUN  CGO_ENABLED=0 GOOS=linux go build -o ratelimiter .


FROM alpine:3.16 
 
ENV REDIS_ADDR=redis:6379
ENV ENVREDIS_USERNAME=""
ENV REDIS_PASSWORD=""
ENV REDIS_DB=0
ENV RL_ADDRESS=":80"
ENV RL_TLL="7d"

WORKDIR /app
COPY --from=build /app/ratelimiter .
CMD ["/app/ratelimiter"]
