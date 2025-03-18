package engine

import (
	"sync"
	"github.com/google/uuid"
)

var sessionStore = make(map[string]int)
var sessionMutex = sync.Mutex{}

// Création d'une session
func CreateSession(userID int) string {
	sessionMutex.Lock()
	defer sessionMutex.Unlock()

	sessionID := uuid.New().String()
	sessionStore[sessionID] = userID
	return sessionID
}

// Vérification d'une session
func GetUserID(sessionID string) (int, bool) {
	sessionMutex.Lock()
	defer sessionMutex.Unlock()

	userID, exists := sessionStore[sessionID]
	return userID, exists
}