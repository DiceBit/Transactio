#(auth server)
echo "Protobuf for auth-service"
cd "/c/Users/danii/PROGRAMMING/GolandProjects/Transactio"
protoc -I . \
      -I$GOPATH/google-api/googleapis \
      internal/auth-service/gRPC/auth.proto \
      --go_out=internal/auth-service/gRPC \
      --go-grpc_out=internal/auth-service/gRPC


#(api api-gateway)
echo "Protobuf for api-gateway"
cd "/c/Users/danii/PROGRAMMING/GolandProjects/Transactio"
protoc -I . \
      -I$GOPATH/google-api/googleapis \
      --grpc-gateway_out=logtostderr=true:internal/api-gateway/gRPC \
      internal/auth-service/gRPC/auth.proto \
      --go_out=internal/api-gateway/gRPC \
      --go-grpc_out=internal/api-gateway/gRPC