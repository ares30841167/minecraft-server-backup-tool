package core

type FileManager interface {
	ScanExistFilesInFolder() error
	TryToSaveFile(fileName string) error
}
