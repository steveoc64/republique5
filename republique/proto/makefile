all:
	protoc -I/usr/local/include -I. \
  	-I${GOPATH}/src \
  	-I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
  	-I${GOPATH}/src/github.com/protocolbuffers/protobuf/src \
  	--go_out=plugins=grpc:. \
  	--grpc-gateway_out=logtostderr=true:. \
  	--grpc-gateway_out=logtostderr=true,grpc_api_configuration=rest.yaml:. \
  	republique.proto
