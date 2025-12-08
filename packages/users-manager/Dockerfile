FROM golang:latest AS build
WORKDIR /app
COPY . .
RUN go build

FROM debian:latest AS base
COPY --from=build /app/users-manager /usr/local/bin/users-manager

CMD [ "users-manager" ]
