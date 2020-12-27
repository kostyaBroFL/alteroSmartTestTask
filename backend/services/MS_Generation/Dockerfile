FROM golang:1.15.6-buster

# cp to /usr/local/go/src/ateroSmartTestTask

ENV PATH $PATH:/gopath/bin
ENV GOPATH /gopath

RUN mkdir -p /protoc && \
    cd /protoc  && \
    apt update -y  && \
    apt install wget autoconf automake libtool curl make g++ unzip tar  && \
    wget https://github.com/protocolbuffers/protobuf/releases/download/v3.14.0/protobuf-all-3.14.0.tar.gz  && \
    tar -xvzf protobuf-all-3.14.0.tar.gz  && \
    cd protobuf-3.14.0  && \
    ./autogen.sh  && \
    ./configure  && \
    make  && \
    make install  && \
    ldconfig && \
    mkdir -p /gopath/bin && \
    mkdir -p /gopath/pkg && \
    mkdir -p /gopath/src
    echo "Installing protobuf" && \
    go get github.com/golang/protobuf/protoc-gen-go@v1.4.3 && \
    go install github.com/golang/protobuf/protoc-gen-go && \
    echo "Installing grpc" && \
    go get google.golang.org/grpc@v1.34.0 && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc && \
    echo "Installing gRPC->RESTful API middleware" && \
    go get github.com/grpc-ecosystem/grpc-gateway@v1 && \
    go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway && \
    go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger && \
    echo "Installing mwitkow/go-proto-validators" && \
    go get -u github.com/mwitkow/go-proto-validators/protoc-gen-govalidators@v0.3.2

ENTRYPOINT ["bash"]
