# Building the binary of the App
FROM golang:1.23 AS build

# boilerplate should be replaced with your project name
WORKDIR /go/src

# Copy all the Code and stuff to compile everything
COPY . .

# Downloads all the dependencies in advance (could be left out, but it's more clear this way)
RUN go mod download

# Builds the application as a staticly linked one, to allow it to run on alpine
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o app .
# RUN CGO_ENABLED=1 GOOS=linux go build -o app -a -ldflags '-linkmode external -extldflags "-static"' .




# Moving the binary to the 'final Image' to make it smaller
FROM alpine:latest
# FROM gcr.io/distroless/base-debian10
# FROM debian:buster-slim
# FROM scratch
RUN apk add --no-cache tzdata
# RUN apk add build-base

WORKDIR /app

# Create the public dir and copy all the assets into it
# RUN mkdir ./static
# COPY ./static ./static

# boilerplate should be replaced here as well
COPY --from=build /go/src/app .

# RUN ls -al
# Exposes port 3000 because our program listens on that port
EXPOSE 3000

CMD ["./app"]
