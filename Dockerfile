FROM golang:latest AS compiling
RUN mkdir -p /go/src/comments
WORKDIR /go/src/comments
ADD . .
WORKDIR /go/src/comments/cmd/server
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM alpine:latest
LABEL version="1.0.0"
LABEL maintainer="Ivan Ivanov<test@test.ru>"
WORKDIR /root/
COPY --from=compiling /go/src/comments/cmd/server/app .
ARG dbhost=192.168.1.35:5432/comments
ENV dbhost="${dbhost}"
ARG dbpass=to_be_redefined_at_conrainer_start
ENV dbpass="${dbpass}"
#ENTRYPOINT ./website
CMD ["./app"]
EXPOSE 8080