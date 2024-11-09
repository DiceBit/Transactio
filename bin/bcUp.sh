#!/bin/bash

echo "Starting bc"
cd $GOPATH/Transactio/internal/blockchain
docker-compose up -d
