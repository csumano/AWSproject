#FROM golang:latest AS build
FROM golang:latest AS build

# Copy source
WORKDIR /home/christian/go/src/htmlServer
COPY . .

# Get required modules (assumes packages have been added to ./vendor)
RUN go get -d -v ./...

# Build a statically-linked Go binary for Linux
RUN CGO_ENABLED=0 GOOS=linux go build -a -o main .

# New build phase -- create binary-only image
FROM alpine:latest

# Add support for HTTPS and time zones
RUN apk update && \
    apk upgrade

# Add support for HTTPS and time zones
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

WORKDIR /root/

# Copy files from previous build container
COPY --from=build /home/christian/go/src/htmlServer/main ./

# Add environment variables
ENV LOGGLY_API_KEY <key>
ENV AWS_ACCESS_KEY_ID <key
ENV AWS_SECRET_ACCESS_KEY <key>
# Check results
RUN env && pwd && find .

# Start the application
CMD ["./main"]