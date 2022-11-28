FROM golang:1.18-alpine as builder
WORKDIR /usr/src/app
COPY ./go.mod ./go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./build/app ./app

FROM alpine:3.17
WORKDIR /usr/src/app
COPY --from=builder /usr/src/app/etc etc
COPY --from=builder /usr/src/app/build/app app
COPY --from=builder /usr/src/app/config.yml config.yml

CMD ["./app"]