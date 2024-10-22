package main

import (
	"Transactio/internal/blockchain/server"
	"context"
)

func main() {
	//TODO: Stop server logic

	srv := server.NewBcSrv()
	srv.RunServer()

	srv.AddBlock(context.Background(),
		"QmdfTbBqBPQ7VNxZEYEj14VmRuZBkqFbiwReogJgS1zR1n", "ownerTest", "fileName", 1024, 1, true)
	srv.AddBlock(context.Background(),
		"Qegihekg...", "ownerTest2", "fileName2", 2048, 1, true)
}
