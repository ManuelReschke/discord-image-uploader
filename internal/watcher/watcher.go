package watcher

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

type Watcher struct {
	fsWatcher         *fsnotify.Watcher
	watchPath         string
	supportedFormats  []string
	deleteAfterUpload bool
	eventChan         chan string
	doneChan          chan bool
	processedFiles    map[string]time.Time
	mutex             sync.RWMutex
}

func New(watchPath string, supportedFormats []string, deleteAfterUpload bool) (*Watcher, error) {
	fsWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("failed to create file system watcher: %w", err)
	}

	if _, err := os.Stat(watchPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("watch path does not exist: %s", watchPath)
	}

	err = fsWatcher.Add(watchPath)
	if err != nil {
		fsWatcher.Close()
		return nil, fmt.Errorf("failed to add watch path: %w", err)
	}

	return &Watcher{
		fsWatcher:         fsWatcher,
		watchPath:         watchPath,
		supportedFormats:  supportedFormats,
		deleteAfterUpload: deleteAfterUpload,
		eventChan:         make(chan string, 100),
		doneChan:          make(chan bool),
		processedFiles:    make(map[string]time.Time),
	}, nil
}

func (w *Watcher) Start() {
	log.Printf("Starting file watcher for path: %s", w.watchPath)

	go w.watchLoop()
	go w.cleanupLoop()
}

func (w *Watcher) Stop() {
	log.Println("Stopping file watcher...")
	close(w.doneChan)
	w.fsWatcher.Close()
	close(w.eventChan)
}

func (w *Watcher) GetEventChan() <-chan string {
	return w.eventChan
}

func (w *Watcher) watchLoop() {
	for {
		select {
		case event, ok := <-w.fsWatcher.Events:
			if !ok {
				return
			}

			if event.Op&fsnotify.Create == fsnotify.Create || event.Op&fsnotify.Write == fsnotify.Write {
				if w.isImageFile(event.Name) && w.shouldProcessFile(event.Name) {
					time.Sleep(100 * time.Millisecond)

					if w.isFileReady(event.Name) {
						w.markFileProcessed(event.Name)
						log.Printf("New image detected: %s", event.Name)
						select {
						case w.eventChan <- event.Name:
						case <-w.doneChan:
							return
						}
					}
				}
			}

		case err, ok := <-w.fsWatcher.Errors:
			if !ok {
				return
			}
			log.Printf("File watcher error: %v", err)

		case <-w.doneChan:
			return
		}
	}
}

func (w *Watcher) isImageFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	for _, supportedExt := range w.supportedFormats {
		if ext == supportedExt {
			return true
		}
	}
	return false
}

func (w *Watcher) isFileReady(filename string) bool {
	stat1, err := os.Stat(filename)
	if err != nil {
		return false
	}

	time.Sleep(50 * time.Millisecond)

	stat2, err := os.Stat(filename)
	if err != nil {
		return false
	}

	return stat1.Size() == stat2.Size() && stat1.ModTime().Equal(stat2.ModTime())
}

func (w *Watcher) DeleteFile(filename string) error {
	if !w.deleteAfterUpload {
		return nil
	}

	err := os.Remove(filename)
	if err != nil {
		return fmt.Errorf("failed to delete file %s: %w", filename, err)
	}

	log.Printf("Deleted file after upload: %s", filename)
	return nil
}

func (w *Watcher) ScanExistingFiles() ([]string, error) {
	var files []string

	err := filepath.Walk(w.watchPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && w.isImageFile(path) {
			files = append(files, path)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to scan existing files: %w", err)
	}

	log.Printf("Found %d existing image files", len(files))
	return files, nil
}

func (w *Watcher) shouldProcessFile(filename string) bool {
	w.mutex.RLock()
	defer w.mutex.RUnlock()

	lastProcessed, exists := w.processedFiles[filename]
	if !exists {
		return true
	}

	return time.Since(lastProcessed) > 5*time.Second
}

func (w *Watcher) markFileProcessed(filename string) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	w.processedFiles[filename] = time.Now()
}

func (w *Watcher) cleanupLoop() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			w.cleanupOldEntries()
		case <-w.doneChan:
			return
		}
	}
}

func (w *Watcher) cleanupOldEntries() {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	cutoff := time.Now().Add(-5 * time.Minute)
	for filename, timestamp := range w.processedFiles {
		if timestamp.Before(cutoff) {
			delete(w.processedFiles, filename)
		}
	}
}
