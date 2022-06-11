package core

import (
	"errors"
	"fmt"
	"log"
	"os"

	"guanyu.dev/minecraft-server-backup-tool/internal/storage"
)

type BackupManager struct {
	storage  storage.Storage
	basePath string
}

func NewBackupManager() (*BackupManager, error) {
	if os.Getenv("WATCH_PATH") == "" {
		return nil, errors.New("BackupManager: 未設定環境變數WATCH_PATH")
	}

	s3, err := storage.NewS3Service()
	if err != nil {
		return nil, err
	}

	return &BackupManager{
		storage:  s3,
		basePath: os.Getenv("WATCH_PATH"),
	}, nil
}

func (bw *BackupManager) TryToSaveFile(fileName string) error {
	log.Println("偵測到新的地圖備份檔:", fileName)

	exist, err := bw.storage.CheckFileIsExist(fileName)
	if err != nil {
		return err
	}

	if exist {
		return errors.New(fmt.Sprintf("地圖備份檔 %s 已經存在於儲存空間中", fileName))
	}

	err = bw.storage.PutFile(fileName)
	if err != nil {
		return err
	}

	log.Printf("地圖備份檔 %s 成功儲存", fileName)
	return nil
}
