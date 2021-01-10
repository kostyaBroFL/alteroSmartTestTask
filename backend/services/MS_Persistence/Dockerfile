FROM golang:1.15.6-buster
ENV PATH $PATH:/go/bin
RUN mkdir -p /protoc && \
    cd /protoc  && \
    apt update -y  && \
    apt install wget autoconf automake libtool curl make g++ unzip tar -y && \
    wget https://github.com/protocolbuffers/protobuf/releases/download/v3.14.0/protoc-3.14.0-linux-x86_64.zip  && \
    unzip protoc-3.14.0-linux-x86_64 -d /usr/local
# gopath /go
WORKDIR /usr/local/go/src/alteroSmartTestTask
COPY . .
RUN go mod download && \
    go install github.com/golang/protobuf/protoc-gen-go && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc && \
    go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway && \
    go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger && \
    go install github.com/mwitkow/go-proto-validators/protoc-gen-govalidators && \
    ./build.sh && \
    go build -o ms_persistence_start backend/services/MS_Persistence/cmd/main.go && \
    chmod +x ms_persistence_start

ENTRYPOINT ["./ms_persistence_start"]

# todo multistage
