FROM golang:alpine as build_container
WORKDIR /app
COPY ./tweet_service/go.mod ./tweet_service/go.sum ./
RUN go mod download
COPY ./tweet_service/ .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest
WORKDIR /usr/bin
COPY --from=build_container /app/main .
EXPOSE 8001
CMD ["./main"]
