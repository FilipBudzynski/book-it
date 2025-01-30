package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/FilipBudzynski/book_it/internal/errs"
	"github.com/FilipBudzynski/book_it/utils"
	"github.com/labstack/echo/v4"
)

type NotificationManager struct {
	clients map[string]chan string
	mu      sync.Mutex
}

func NewConnectionManager() *NotificationManager {
	return &NotificationManager{
		clients: make(map[string]chan string),
	}
}

func (cm *NotificationManager) AddClient(userID string, ch chan string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	if _, exists := cm.clients[userID]; !exists {
		cm.clients[userID] = ch
		fmt.Println("Added client:", userID)
	}
}

func (cm *NotificationManager) RemoveClient(userID string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	if ch, exists := cm.clients[userID]; exists {
		close(ch)
		delete(cm.clients, userID)
		fmt.Println("Removed client:", userID)
	}
}

func (cm *NotificationManager) GetClientChannel(userID string) (chan string, bool) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	ch, ok := cm.clients[userID]
	return ch, ok
}

func (cm *NotificationManager) SseHandler(c echo.Context) error {
	w := c.Response()
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	userID, err := utils.GetUserEmailFromSession(c.Request())
	if err != nil {
		return errs.HttpErrorUnauthorized(err)
	}
	dataCh, ok := cm.GetClientChannel(userID)
	if !ok {
		dataCh = make(chan string, 10)
		cm.AddClient(userID, dataCh)
	}
	ctx, cancel := context.WithCancel(c.Request().Context())
	defer func() {
		cancel()
		cm.RemoveClient(userID)
	}()
	go func() {
		for {
			select {
			case data, ok := <-dataCh:
				if !ok {
					return
				}
				_, _ = fmt.Fprintf(w.Writer, "data: %s\n\n", data)
				if flusher, ok := w.Writer.(http.Flusher); ok {
					flusher.Flush()
				}
			case <-ctx.Done():
				return
			}
		}
	}()
	<-ctx.Done()
	log.Println("Client disconnected:", userID)
	return nil
}

func (cm *NotificationManager) Notify(userID string, message string) {
	if msgChannel, ok := cm.GetClientChannel(userID); ok {
		select {
		case msgChannel <- message:
		default:
			log.Println("Channel full or closed, message not sent")
		}
	} else {
		log.Println("No active channel for the user")
	}
}
