package utils

import (
	"bytes"
	"io"
	"mime/multipart"
)

func ConvertMultipartToBytes(file multipart.File) ([]byte, error) {
	defer file.Close()

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
