package core

type FileManager interface {
	TryToSaveFile(fileName string) error
}
