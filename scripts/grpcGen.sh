#(auth server)
echo "Protobuf for auth-service"
cd "/c/Users/danii/PROGRAMMING/GolandProjects/Transactio"
protoc -I . \
      -I$GOPATH/google-api/googleapis \
      internal/auth-service/gRPC/auth.proto \
      --go_out=internal/auth-service/gRPC \
      --go-grpc_out=internal/auth-service/gRPC


#(api gateway)
echo "Protobuf for api-gateway"
cd "/c/Users/danii/PROGRAMMING/GolandProjects/Transactio"
protoc -I . \
      -I$GOPATH/google-api/googleapis \
      --grpc-gateway_out=logtostderr=true:internal/gateway/gRPC \
      internal/auth-service/gRPC/auth.proto \
      --go_out=internal/gateway/gRPC \
      --go-grpc_out=internal/gateway/gRPC