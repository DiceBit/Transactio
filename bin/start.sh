#!/bin/bash
cd "$GOPATH/Transactio/cmd"

echo "Starting gateway"
cd "api-gateway/"
go run ./main.go &
PID_GATEWAY=$!

echo "Starting auth-service"
cd "../auth-service/"
./main.exe &
PID_AUTH=$!

# Ожидание завершения всех процессов
wait $PID_GATEWAY $PID_AUTH
echo "All services completed!"


