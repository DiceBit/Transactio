package models

import (
	"strconv"
)

// File metadata
type FileMD struct {
	Cid       string
	OwnerAddr string // Адрес владельца файла

	FileName string
	FileSize int

	IsDelete  bool // active or delete
	IsSecured bool
}

func NewData(
	cid, owner, fileName string,
	fileSize int, isDelete, isSecured bool) *FileMD {

	data := &FileMD{
		Cid:       cid,
		OwnerAddr: owner,

		FileName: fileName,
		FileSize: fileSize,

		IsDelete:  isDelete,
		IsSecured: isSecured,
	}
	return data
}

func (f *FileMD) fmdToString() string {
	header := f.Cid + f.OwnerAddr + f.FileName +
		strconv.Itoa(f.FileSize) + strconv.FormatBool(f.IsDelete) + strconv.FormatBool(f.IsSecured)
	return header
}
