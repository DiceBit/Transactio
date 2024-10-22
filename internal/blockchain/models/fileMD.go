package models

import (
	"strconv"
	"time"
)

// File metadata
type FileMD struct {
	Cid       string
	OwnerAddr string // Адрес владельца файла

	FileName string
	FileSize int
	CreateAt int64

	Status  bool // active or delete
	Version int  // Версия файла
}

func NewData(
	cid, owner, fileName string,
	fileSize, version int,
	status bool) *FileMD {

	timestamp := time.Now().Unix()
	data := &FileMD{
		Cid:       cid,
		OwnerAddr: owner,

		FileName: fileName,
		FileSize: fileSize,
		CreateAt: timestamp,

		Status:  status,
		Version: version,
		//accessList: accessList,
	}
	return data
}

func (f *FileMD) fmdToString() string {
	timestamp := strconv.FormatInt(f.CreateAt, 10)
	header := f.Cid + f.OwnerAddr + f.FileName +
		strconv.Itoa(f.FileSize) + timestamp + strconv.FormatBool(f.Status) + strconv.Itoa(f.Version)
	return header
}
