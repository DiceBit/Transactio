#!/usr/bin/env bash

# Пути к директориям
fs_dir="$GOPATH/Transactio/internal/fileStorage"
bc_dir="$GOPATH/Transactio/internal/blockchain"
fs_cmd_dir="$GOPATH/Transactio/cmd/fsService"
bc_cmd_dir="$GOPATH/Transactio/cmd/blockchain"


export ipfs_staging="$GOPATH/Transactio/pkg/ipfsStorage/staging"
export ipfs_data="$GOPATH/Transactio/pkg/ipfsStorage/data"

# Запуск сервисов с объединением вывода
echo "Starting fs"
cd "$fs_dir"
docker-compose up -d &> /dev/tty &
cd "$fs_cmd_dir"
go run ./main.go &> /dev/tty &
PID_FS=$!

echo "Starting bc"
cd "$bc_dir"
docker-compose up -d &> /dev/tty &
cd "$bc_cmd_dir"
go run ./main.go &> /dev/tty &
PID_BC=$!

# Ожидание завершения сервисов
wait $PID_FS $PID_BC
echo "All services completed!"