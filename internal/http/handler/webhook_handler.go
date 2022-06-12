package handler

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"guanyu.dev/minecraft-server-backup-tool/internal/core"
	"guanyu.dev/minecraft-server-backup-tool/internal/http/model"
)

type WebhookServerHandlers struct {
	bakmanager core.FileManager
}

func NewWebhookServerHandlers() (*WebhookServerHandlers, error) {
	manager, err := core.NewBackupManager()
	if err != nil {
		return nil, err
	}

	return &WebhookServerHandlers{
		bakmanager: manager,
	}, nil
}

func (h *WebhookServerHandlers) WebhookHandler(w http.ResponseWriter, req *http.Request) {
	if req.Header.Get("content-type") != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var backupDoneNotification model.BackupDoneNotification
	err = json.Unmarshal(body, &backupDoneNotification)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if backupDoneNotification.Event == nil ||
		backupDoneNotification.BackupFilename == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if *backupDoneNotification.Event != "DONE" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.bakmanager.TryToSaveFile(*backupDoneNotification.BackupFilename)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "")
}
