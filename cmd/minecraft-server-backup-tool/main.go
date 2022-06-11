package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/joho/godotenv/autoload"
	"guanyu.dev/minecraft-server-backup-tool/internal/core"
)

func main() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)

	bakupManager, err := core.NewBackupManager()
	if err != nil {
		log.Fatal(err.Error())
	}

	err = bakupManager.ScanExistFilesInFolder()
	if err != nil {
		log.Fatal(err.Error())
	}

	bakupWatcher, err := core.NewBackupWatcher()
	if err != nil {
		log.Fatal(err.Error())
	}

	bakupWatcher.StartWatchFile()

	log.Println("開始監聽地圖備份資料夾")
	<-sig
}
