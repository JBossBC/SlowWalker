package utils

import (
	"errors"
	"io"
	"net/http"
	"os"
)

// the zero copy to download file
func WriteFile(stream http.ResponseWriter, fileName string) {
	_, err := os.Stat(fileName)
	if errors.Is(err, &os.PathError{}) {
		stream.Write(NewFailedResponse("下载失败").Serialize())
		return
	}
	src, _ := os.Open(fileName)
	io.Copy(stream, src)
}
