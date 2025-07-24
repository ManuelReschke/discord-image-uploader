package uploader

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"discord-image-uploader/internal/config"
	"discord-image-uploader/internal/discord"
	"discord-image-uploader/internal/history"
	"discord-image-uploader/internal/watcher"
)

type Uploader struct {
	config        *config.Config
	discordClient *discord.Client
	watcher       *watcher.Watcher
	history       *history.History
	queue         []string
	queueMutex    sync.RWMutex
	ticker        *time.Ticker
	doneChan      chan bool
}

func New(cfg *config.Config, discordClient *discord.Client, watcher *watcher.Watcher, history *history.History) *Uploader {
	return &Uploader{
		config:        cfg,
		discordClient: discordClient,
		watcher:       watcher,
		history:       history,
		queue:         make([]string, 0),
		doneChan:      make(chan bool),
	}
}

func (u *Uploader) Start() error {
	log.Println("Starting uploader...")

	if u.config.History.CleanupMissingFiles {
		err := u.history.CleanupMissingFiles()
		if err != nil {
			log.Printf("Warning: failed to cleanup missing files from history: %v", err)
		}
	}

	existingFiles, err := u.watcher.ScanExistingFiles()
	if err != nil {
		return fmt.Errorf("failed to scan existing files: %w", err)
	}

	u.addToQueue(existingFiles...)

	u.ticker = time.NewTicker(time.Duration(u.config.Upload.IntervalSeconds) * time.Second)

	go u.processQueue()
	go u.watchForNewFiles()

	return nil
}

func (u *Uploader) Stop() {
	log.Println("Stopping uploader...")

	if u.ticker != nil {
		u.ticker.Stop()
	}

	close(u.doneChan)

	u.processRemainingQueue()
}

func (u *Uploader) addToQueue(files ...string) {
	u.queueMutex.Lock()
	defer u.queueMutex.Unlock()

	var newFiles []string
	for _, file := range files {
		if u.isValidFile(file) && !u.history.IsUploaded(file) {
			u.queue = append(u.queue, file)
			newFiles = append(newFiles, file)
		} else if u.history.IsUploaded(file) {
			log.Printf("Skipping already uploaded file: %s", file)
		}
	}

	if len(newFiles) > 0 {
		log.Printf("Added %d new files to queue", len(newFiles))
	}
}

func (u *Uploader) processQueue() {
	for {
		select {
		case <-u.ticker.C:
			u.uploadBatch()
		case <-u.doneChan:
			return
		}
	}
}

func (u *Uploader) processRemainingQueue() {
	u.queueMutex.RLock()
	queueLength := len(u.queue)
	u.queueMutex.RUnlock()

	if queueLength > 0 {
		log.Printf("Processing remaining %d files in queue...", queueLength)
		u.uploadBatch()
	}
}

func (u *Uploader) uploadBatch() {
	u.queueMutex.Lock()
	defer u.queueMutex.Unlock()

	if len(u.queue) == 0 {
		return
	}

	batchSize := u.config.Upload.BatchSize
	if batchSize > len(u.queue) {
		batchSize = len(u.queue)
	}

	batch := u.queue[:batchSize]
	u.queue = u.queue[batchSize:]

	log.Printf("Uploading batch of %d files", len(batch))

	if len(batch) == 1 {
		err := u.discordClient.UploadImage(batch[0])
		if err != nil {
			log.Printf("Failed to upload %s: %v", batch[0], err)
			u.queue = append([]string{batch[0]}, u.queue...)
			return
		}
		u.handleSuccessfulUpload(batch[0])
	} else {
		err := u.discordClient.UploadImages(batch)
		if err != nil {
			log.Printf("Failed to upload batch: %v", err)
			u.queue = append(batch, u.queue...)
			return
		}

		for _, file := range batch {
			u.handleSuccessfulUpload(file)
		}
	}
}

func (u *Uploader) watchForNewFiles() {
	eventChan := u.watcher.GetEventChan()

	for {
		select {
		case file, ok := <-eventChan:
			if !ok {
				return
			}
			u.addToQueue(file)

		case <-u.doneChan:
			return
		}
	}
}

func (u *Uploader) handleSuccessfulUpload(file string) {
	err := u.history.MarkUploaded(file, "")
	if err != nil {
		log.Printf("Warning: failed to mark file as uploaded in history: %v", err)
	}

	err = u.watcher.DeleteFile(file)
	if err != nil {
		log.Printf("Warning: failed to delete file after upload: %v", err)
	}
}

func (u *Uploader) isValidFile(file string) bool {
	stat, err := os.Stat(file)
	if err != nil {
		log.Printf("Cannot stat file %s: %v", file, err)
		return false
	}

	maxSizeBytes := int64(u.config.Upload.MaxFileSizeMB) * 1024 * 1024
	if stat.Size() > maxSizeBytes {
		log.Printf("File %s is too large (%d bytes, max: %d bytes)", file, stat.Size(), maxSizeBytes)
		return false
	}

	return true
}

func (u *Uploader) GetQueueLength() int {
	u.queueMutex.RLock()
	defer u.queueMutex.RUnlock()
	return len(u.queue)
}
