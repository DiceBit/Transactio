#(api api-gateway)
echo "Protobuf for api-gateway"
cd "$GOPATH/Transactio/"
protoc -I . \
      -I$GOPATH/google-api/googleapis \
      --grpc-gateway_out=logtostderr=true:internal/api-gateway/gRPC \
      internal/auth-service/gRPC/auth.proto \
      --go_out=internal/api-gateway/gRPC \
      --go-grpc_out=internal/api-gateway/gRPC
echo "Completed!"


#(auth server)
echo "Protobuf for auth-service"
protoc -I . \
      -I$GOPATH/google-api/googleapis \
      internal/auth-service/gRPC/auth.proto \
      --go_out=internal/auth-service/gRPC \
      --go-grpc_out=internal/auth-service/gRPC
echo "Completed!"


#(blockchain server)
echo "Protobuf for blockchain"
protoc -I . \
      internal/fileStorage/gRPC/fileStorage.proto \
      --go_out=internal/blockchain/gRPC \
      --go-grpc_out=internal/blockchain/gRPC
echo "Completed!"


#(fileStorage server)
echo "Protobuf for fileService"
protoc -I . \
      internal/fileStorage/gRPC/fileStorage.proto \
      --go_out=internal/fileStorage/gRPC \
      --go-grpc_out=internal/fileStorage/gRPC
      echo "Completed!"