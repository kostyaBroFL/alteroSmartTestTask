protoc --go_out=plugins=grpc:. \
    --proto_path=${GOPATH}/pkg/mod/github.com/mwitkow/go-proto-validators\@v0.3.2 \
    --proto_path=${GOPATH}/pkg/mod/github.com/grpc-ecosystem/grpc-gateway\@v1.16.0/third_party/googleapis \
    --proto_path=. \
    --grpc-gateway_out ./ \
    --grpc-gateway_opt logtostderr=true \
    --grpc-gateway_opt paths=source_relative \
    --grpc-gateway_opt generate_unbound_methods=true \
    --govalidators_out=. \
    service.proto
