package service

import (
	"sync"
)

type FileService struct {
}

var (
	fileService *FileService
	fileOnce    sync.Once
)

func getFileService() *FileService {
	fileOnce.Do(func() {
		fileService = &FileService{}
	})
	return fileService
}
func (file *FileService) AuthorizeFiles(files []string, user string) {

}
