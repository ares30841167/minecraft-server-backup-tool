package http

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"guanyu.dev/minecraft-server-backup-tool/internal/http/handler"
)

type WebhookServer struct {
	server *http.Server
}

func TcpPortIsOccupied(port int) bool {
	timeout := 5 * time.Second
	target := fmt.Sprintf(":%d", port)

	conn, err := net.DialTimeout("tcp", target, timeout)
	if err != nil {
		return false
	}

	if conn != nil {
		conn.Close()
		log.Println(fmt.Sprintf("WebhookServer: 連接埠%d已被占用", port))
		return true
	}

	log.Println(fmt.Sprintf("WebhookServer: 連接埠%d未被占用", port))
	return false
}

func ChooseRandomPort() int {
	newPort := 6667 + rand.Intn(100)
	log.Println(fmt.Sprintf("WebhookServer: 隨機選擇新連接埠: %d", newPort))
	return newPort
}

func GetAvaliablePort() (int, error) {
	counter := 0
	reservePort := 6666

	for {
		if counter == 10 {
			return -1, errors.New("隨機選取連接埠動作次數已達10次上限!")
		}

		if TcpPortIsOccupied(reservePort) {
			reservePort = ChooseRandomPort()
			counter = counter + 1
		} else {
			log.Println(fmt.Sprintf("WebhookServer: 將使用連接埠%d啟動Webhook伺服器", reservePort))
			break
		}
	}

	return reservePort, nil
}

func NewWebhookServer() (*WebhookServer, error) {

	port, err := GetAvaliablePort()
	if err != nil {
		return nil, err
	}

	mux := http.NewServeMux()

	handlers, err := handler.NewWebhookServerHandlers()
	if err != nil {
		return nil, err
	}

	mux.HandleFunc("/", handlers.WebhookHandler)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      mux,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	}

	return &WebhookServer{
		server: server,
	}, nil
}

func (ws *WebhookServer) Start() error {
	log.Println("WebhookServer: 開始監聽Webhook通知")
	return ws.server.ListenAndServe()
}

func (ws *WebhookServer) Close() error {
	return ws.server.Close()
}

func (ws *WebhookServer) StartInterruptListener() chan struct{} {
	processed := make(chan struct{})

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		if err := ws.server.Shutdown(ctx); nil != err {
			log.Fatalf("WebhookServer: 伺服器關閉失敗, 錯誤: %v\n", err)
		}
		log.Println("WebhookServer: 已關閉Webhook伺服器")
		close(processed)
	}()

	return processed
}
