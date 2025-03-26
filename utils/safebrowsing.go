package utils

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/google/safebrowsing"
)

var (
	safeBrowsingClient *safebrowsing.SafeBrowser
	sbInitOnce         sync.Once
	sbEnabled          bool
)

// initializes the Safe Browsing client
func InitSafeBrowsing(apiKey string) error {
	var initErr error
	sbInitOnce.Do(func() {
		if apiKey == "" {
			log.Println("Safe Browsing API key not provided, skipping initialization")
			return
		}

		config := safebrowsing.Config{
			APIKey: apiKey,
			DBPath: "/tmp/safebrowsing.db",
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		sb, err := safebrowsing.NewSafeBrowser(config)
		if err != nil {
			initErr = err
			return
		}

		// Warm up the DB
		if err := sb.WaitUntilReady(ctx); err != nil {
			initErr = err
			return
		}

		safeBrowsingClient = sb
		sbEnabled = true
		log.Println("Safe Browsing initialized successfully")
	})
	return initErr
}

// checks if a URL is safe using Google Safe Browsing
func CheckURLSafety(url string) (bool, error) {
	if !sbEnabled || safeBrowsingClient == nil {
		return true, nil
	}

	threats, err := safeBrowsingClient.LookupURLs([]string{url})
	if err != nil {
		return false, err
	}

	return len(threats[0]) == 0, nil
}
