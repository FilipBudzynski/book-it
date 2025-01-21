package handlers

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/FilipBudzynski/book_it/internal/errs"
	"github.com/FilipBudzynski/book_it/utils"
	"github.com/labstack/echo/v4"
)

// A struct to manage active SSE connections
type ConnectionManager struct {
	clients map[string]chan string // Channels per user
	mu      sync.Mutex             // To protect the clients map
}

// NewConnectionManager creates a new instance of ConnectionManager
func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		clients: make(map[string]chan string),
	}
}

// AddClient adds a new client and their message channel to the manager
func (cm *ConnectionManager) AddClient(userID string, ch chan string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.clients[userID] = ch
	fmt.Println("Added client: ", userID)
}

// RemoveClient removes a client from the manager
func (cm *ConnectionManager) RemoveClient(userID string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	if ch, exists := cm.clients[userID]; exists {
		close(ch) // Close the channel to stop sending messages
		delete(cm.clients, userID)
		fmt.Println("Removed client:", userID)
	}
}

// GetClientChannel returns the message channel for a user
func (cm *ConnectionManager) GetClientChannel(userID string) (chan string, bool) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	ch, ok := cm.clients[userID]
	return ch, ok
}

func (cm *ConnectionManager) SseHandler(c echo.Context) error {
	// Set headers for SSE
	w := c.Response()
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	userID, err := utils.GetUserEmailFromSession(c.Request())
	if err != nil {
		return errs.HttpErrorUnauthorized(err)
	}

	// Get or create the user's message channel
	dataCh, ok := cm.GetClientChannel(userID)
	if !ok {
		dataCh = make(chan string)
		cm.AddClient(userID, dataCh)

	}

	// Create a context for handling client disconnection
	_, cancel := context.WithCancel(c.Request().Context())
	defer cancel()

	// Send data to the client
	for {
		select {
		case data := <-dataCh:
			// Send the message to the client
			fmt.Fprintf(c.Response().Writer, "data: %s\n\n", data)

			// Ensure the data is flushed to the client
			if flusher, ok := c.Response().Writer.(http.Flusher); ok {
				flusher.Flush()
			}

		case <-c.Request().Context().Done():
			// This will trigger when the client disconnects
			fmt.Println("Client disconnected")
			cm.RemoveClient(userID)
			return nil
		}
	}
}
