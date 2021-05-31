protoc -I protos --go_out=plugins=grpc:. --grpc-gateway_out=. protos/helloworld/*.proto
protoc -I protos --go_out=plugins=grpc:. --grpc-gateway_out=. protos/health/*.proto