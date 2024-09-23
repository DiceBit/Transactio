#(user server)
echo "Protobuf for auth-service"
cd "/c/Users/danii/PROGRAMMING/GolandProjects/Ecommers"
protoc -I . \
      -I$GOPATH/google-api/googleapis \
      services/auth-service/pkg/gRPC/auth.proto \
      --go_out=services/auth-service/pkg/gRPC \
      --go-grpc_out=services/auth-service/pkg/gRPC

#ДЛЯ БУДУЩЕГО
#т.к Post хочет взаимодействовать с user-server'ом, то ему необходимо создать client
#схема https://habrastorage.org/r/w1560/webt/sh/_w/qg/sh_wqgsdmwcbhrfcvv0f5sluwo4.png

#(client для user сервиса)
#protoc -I. \
#   -I$GOPATH/Ecommers/google-api/googleapis \
#   services/auth-service/gRPC/auth.proto --go_out=plugins=grpc:./services/post/protobuf/

#(api gateway)
echo "Protobuf for api-gateway"
cd "/c/Users/danii/PROGRAMMING/GolandProjects/Ecommers"
protoc -I . \
      -I$GOPATH/google-api/googleapis \
      --grpc-gateway_out=logtostderr=true:services/api-gateway/gRPC \
      services/auth-service/pkg/gRPC/auth.proto \
      --go_out=services/api-gateway/gRPC \
      --go-grpc_out=services/api-gateway/gRPC