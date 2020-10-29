FROM golang:1.15.3 as build

# Install the Protocol Buffers compiler and Go plugin
RUN apt-get update && apt-get install -y zip
RUN go get \
    github.com/golang/protobuf/protoc-gen-go \
    google.golang.org/grpc/cmd/protoc-gen-go-grpc

# Create the source folder
WORKDIR /go/plugin

# Copy the source to the build folder
COPY . .

# Build and zip the plugin
RUN make all

FROM scratch as export

COPY --from=build /go/plugin/bin/*.zip .