#docker build \
#  --tag 'ms_generation:latest' \
#  --file backend/services/MS_Generation/Dockerfile .
cd ./backend/services/MS_Generation/common/api
./build.sh
cd ../../../../..

cd ./backend/services/MS_Persistence/common/api
./build.sh
cd ../../../../..

#protoc --go_out=plugins=grpc:. \
#    --proto_path=${GOPATH}/src/github.com/mwitkow/go-proto-validators\@v0.3.2 \
#    --proto_path=${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway\@v1.16.0/third_party/googleapis \
#    --proto_path=. \
#    --govalidators_out=. \
#    backend/services/MS_Generation/common/api/common/device.proto && \
#protoc --go_out=plugins=grpc:. \
#    --proto_path=${GOPATH}/src/github.com/mwitkow/go-proto-validators\@v0.3.2 \
#    --proto_path=${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway\@v1.16.0/third_party/googleapis \
#    --grpc-gateway_out ./ \
#    --grpc-gateway_opt logtostderr=true \
#    --grpc-gateway_opt paths=source_relative \
#    --grpc-gateway_opt generate_unbound_methods=true \
#    --proto_path=. \
#    --govalidators_out=. \
#    backend/services/MS_Generation/common/api/service.proto