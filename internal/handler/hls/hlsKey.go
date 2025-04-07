package hls

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

const keyDir = "./keys/"

var (
	keyCache   = make(map[string][]byte)
	cacheMutex sync.RWMutex
)

// InitKeys is a function that initializes all keys
// @Summary Initialize all keys
// @Description Initialize all keys
// @Tags Hls
func InitKeys() {

	if err := os.MkdirAll(keyDir, 0755); err != nil {
		log.Fatalf("Failed to create key directory: %v", err)
	}

	files, err := os.ReadDir(keyDir)
	if err != nil {
		log.Fatalf("Failed to read key directory: %v", err)
	}

	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".key") {
			keyPath := filepath.Join(keyDir, file.Name())

			keyData, err := os.ReadFile(keyPath)
			if err != nil {
				log.Printf("Failed to read key file: %s (%v)", keyPath, err)
				continue
			}

			// store key in cache
			keyCache[file.Name()] = keyData
			log.Printf("Loading Key: %s (%d bytes)", file.Name(), len(keyData))
		}
	}

	log.Println("All key files have been loaded into memory!")
}

// KeyHandler is a function that handles key requests
// @Summary Key handler
// @Description Key handler
// @Tags Hls
// @Accept  json
func KeyHandler(c *gin.Context) {
	//keyName := c.Query("key")
	keyName := c.PostForm("key")
	if keyName == "" || !strings.HasSuffix(keyName, ".key") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid key file request"})
		return
	}

	keyName = filepath.Clean(keyName)

	// get key form cache
	cacheMutex.RLock()
	key, exists := keyCache[keyName]
	cacheMutex.RUnlock()

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Key file not found"})
		return
	}
	c.Data(http.StatusOK, "application/octet-stream", key)
}
