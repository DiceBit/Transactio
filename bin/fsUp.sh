#!/bin/bash

echo "Starting fs"
export ipfs_staging="$GOPATH/Transactio/pkg/ipfsStorage/staging"
export ipfs_data="$GOPATH/Transactio/pkg/ipfsStorage/data"

cd $GOPATH/Transactio/internal/fileStorage
docker-compose up