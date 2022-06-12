package model

type BackupDoneNotification struct {
	Event          *string `json:"event"`
	BackupFilename *string `json:"backupFilename"`
}
