package models

import (
	"crypto/sha256"
	"fmt"
	"strconv"
	"time"
)

type Blockchain struct {
	Index         int
	Hash          string
	PrevBlockHash string
	Timestamp     int64

	Fmd *FileMD
}

func (bc *Blockchain) SetHash() {
	timestamp := strconv.FormatInt(bc.Timestamp, 10)
	headers := bc.Fmd.fmdToString() + bc.PrevBlockHash + timestamp
	hash := sha256.Sum256([]byte(headers))

	bc.Hash = fmt.Sprintf("%x", hash[:])
}

func SetGenesisBlock() *Blockchain {
	dataGenesis := NewData("", "genesis", "Genesis Block", 0, false, false)
	genesisBlock := Blockchain{
		PrevBlockHash: "",
		Timestamp:     time.Now().Unix(),
		Fmd:           dataGenesis,
	}
	genesisBlock.SetHash()

	return &genesisBlock
}
