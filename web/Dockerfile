FROM golang:alpine as build

ENV GO111MODULE=on

RUN apk update && apk --no-cache add bash ca-certificates git gcc libc-dev

RUN mkdir -p /app/src/handler

WORKDIR /app/src/handler

COPY go.mod .
COPY go.sum .
COPY main.go .

RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o /app/src/handler/bin/todos /app/src/handler/main.go
RUN go test $(go list ./... | grep -v /vendor/) -cover

FROM alpine:3.8

# Add non root user and certs
RUN apk --no-cache add ca-certificates \
  && addgroup -S app && adduser -S -g app app \
  && mkdir -p /home/app/static \
  && chown app /home/app

WORKDIR /home/app

COPY --from=build /app/src/handler/bin/todos .

USER app
EXPOSE 9099

CMD ./todos -d "/home/app/static"
