package service

type File interface {
	AuthorizeFiles(files []string, user string)
}

// var (
// 	fileService *service.FileService
// 	fileOnce    sync.Once
// )

func GetFileService() File {
	return getFileService()
}
