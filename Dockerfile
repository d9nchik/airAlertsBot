#build stage
FROM golang:alpine AS builder
RUN apk add --no-cache git
WORKDIR /go/src/app
RUN mkdir -p /go/bin
COPY . .
RUN go get -d -v ./...
RUN cd cmd && go build -o /go/bin/app

#final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/bin/app /app
ENTRYPOINT /app