package main

import (
	"log"

	_ "github.com/joho/godotenv/autoload"
	"guanyu.dev/minecraft-server-backup-tool/internal/core"
	"guanyu.dev/minecraft-server-backup-tool/internal/http"
)

func main() {
	bakupManager, err := core.NewBackupManager()
	if err != nil {
		log.Fatal(err.Error())
	}

	err = bakupManager.ScanExistFilesInFolder()
	if err != nil {
		log.Fatal(err.Error())
	}

	server, err := http.NewWebhookServer()
	if err != nil {
		log.Fatal(err.Error())
	}

	shutdownProcessed := server.StartInterruptListener()
	server.Start()

	<-shutdownProcessed
	defer server.Close()
}
