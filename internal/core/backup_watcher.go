package core

import (
	"errors"
	"log"
	"os"
	"regexp"

	"github.com/fsnotify/fsnotify"
)

type BackupWatcher struct {
	fswatcher  *fsnotify.Watcher
	bakmanager FileManager
	regexp     *regexp.Regexp
}

func NewBackupWatcher() (*BackupWatcher, error) {
	if os.Getenv("WATCH_PATH") == "" {
		return nil, errors.New("BackupWatcher: 未設定環境變數WATCH_PATH")
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	err = watcher.Add(os.Getenv("WATCH_PATH"))
	if err != nil {
		return nil, err
	}

	manager, err := NewBackupManager()
	if err != nil {
		return nil, err
	}

	r, err := regexp.Compile("^[0-9]+-[0-9]+-[0-9]+-[0-9]+-[0-9]+-[0-9]+\\.zip$")
	if err != nil {
		return nil, err
	}

	return &BackupWatcher{
		fswatcher:  watcher,
		bakmanager: manager,
		regexp:     r,
	}, nil
}

func (bw *BackupWatcher) StartWatchFile() {
	go bw.fsEventListener()
}

func (bw *BackupWatcher) fsEventListener() {
	for {
		select {
		case event, ok := <-bw.fswatcher.Events:
			if !ok {
				return
			}

			if event.Op != fsnotify.Create {
				continue
			}

			if !bw.regexp.MatchString(event.Name) {
				continue
			}

			err := bw.bakmanager.TryToSaveFile(event.Name)
			if err != nil {
				log.Println(err.Error())
			}
		case err, ok := <-bw.fswatcher.Errors:
			if !ok {
				return
			}
			log.Println("發生錯誤:", err)
		}
	}
}
