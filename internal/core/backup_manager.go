package core

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"

	"guanyu.dev/minecraft-server-backup-tool/internal/storage"
)

type BackupManager struct {
	storage  storage.Storage
	basePath string
	regexp   *regexp.Regexp
}

func NewBackupManager() (*BackupManager, error) {
	if os.Getenv("WATCH_PATH") == "" {
		return nil, errors.New("BackupManager: 未設定環境變數WATCH_PATH")
	}

	if os.Getenv("BACKUP_FILE_REGEXP") == "" {
		return nil, errors.New("BackupWatcher: 未設定環境變數BACKUP_FILE_REGEXP")
	}

	r, err := regexp.Compile(os.Getenv("BACKUP_FILE_REGEXP"))
	if err != nil {
		return nil, err
	}

	s3, err := storage.NewS3Service()
	if err != nil {
		return nil, err
	}

	return &BackupManager{
		storage:  s3,
		basePath: os.Getenv("WATCH_PATH"),
		regexp:   r,
	}, nil
}

func (bw *BackupManager) ScanExistFilesInFolder() error {
	log.Println("掃描地圖備份資料夾內現有檔案並嘗試上傳...")

	files, err := ioutil.ReadDir(bw.basePath)
	if err != nil {
		return err
	}

	for _, f := range files {
		if !bw.regexp.MatchString(f.Name()) {
			continue
		}
		err := bw.TryToSaveFile(f.Name())
		if err != nil {
			log.Println(err.Error())
		}
	}

	log.Println("掃描完畢!")
	return nil
}

func (bw *BackupManager) TryToSaveFile(fileName string) error {
	if !bw.regexp.MatchString(fileName) {
		return nil
	}

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
