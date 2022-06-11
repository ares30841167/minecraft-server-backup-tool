package storage

type Storage interface {
	GetFileList() ([]string, error)
	CheckFileIsExist(fileName string) (bool, error)
	PutFile(backupFileName string) error
}
