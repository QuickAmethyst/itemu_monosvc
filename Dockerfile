FROM golang:1.18 as builder
WORKDIR /usr/src/app
COPY . .
RUN make build

FROM golang:1.18
WORKDIR /usr/src/app
COPY --from=builder /usr/src/app/build/app build/app
COPY --from=builder /usr/src/app/etc etc
COPY --from=builder /usr/src/app/config.yml config.yml
EXPOSE 8000
CMD ./build/app