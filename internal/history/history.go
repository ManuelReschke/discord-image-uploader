package history

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type UploadRecord struct {
	FilePath   string    `json:"file_path"`
	FileHash   string    `json:"file_hash"`
	FileSize   int64     `json:"file_size"`
	UploadedAt time.Time `json:"uploaded_at"`
	DiscordURL string    `json:"discord_url,omitempty"`
}

type History struct {
	records     map[string]UploadRecord
	historyFile string
	mutex       sync.RWMutex
}

func New(historyFile string) (*History, error) {
	h := &History{
		records:     make(map[string]UploadRecord),
		historyFile: historyFile,
	}

	err := h.load()
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to load history: %w", err)
	}

	log.Printf("Loaded %d upload records from history", len(h.records))
	return h, nil
}

func (h *History) IsUploaded(filePath string) bool {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return false
	}

	hash, err := h.calculateFileHash(filePath)
	if err != nil {
		return false
	}

	record, exists := h.records[filePath]
	if !exists {
		return false
	}

	return record.FileHash == hash && record.FileSize == fileInfo.Size()
}

func (h *History) MarkUploaded(filePath string, discordURL string) error {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return fmt.Errorf("failed to get file info: %w", err)
	}

	hash, err := h.calculateFileHash(filePath)
	if err != nil {
		return fmt.Errorf("failed to calculate file hash: %w", err)
	}

	record := UploadRecord{
		FilePath:   filePath,
		FileHash:   hash,
		FileSize:   fileInfo.Size(),
		UploadedAt: time.Now(),
		DiscordURL: discordURL,
	}

	h.mutex.Lock()
	h.records[filePath] = record
	h.mutex.Unlock()

	return h.save()
}

func (h *History) RemoveRecord(filePath string) error {
	h.mutex.Lock()
	delete(h.records, filePath)
	h.mutex.Unlock()

	return h.save()
}

func (h *History) GetUploadCount() int {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	return len(h.records)
}

func (h *History) CleanupMissingFiles() error {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	var toRemove []string
	for filePath := range h.records {
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			toRemove = append(toRemove, filePath)
		}
	}

	for _, filePath := range toRemove {
		delete(h.records, filePath)
		log.Printf("Removed missing file from history: %s", filePath)
	}

	if len(toRemove) > 0 {
		return h.save()
	}

	return nil
}

func (h *History) calculateFileHash(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

func (h *History) load() error {
	data, err := os.ReadFile(h.historyFile)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &h.records)
}

func (h *History) save() error {
	dir := filepath.Dir(h.historyFile)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create history directory: %w", err)
	}

	h.mutex.RLock()
	data, err := json.MarshalIndent(h.records, "", "  ")
	h.mutex.RUnlock()

	if err != nil {
		return fmt.Errorf("failed to marshal history: %w", err)
	}

	return os.WriteFile(h.historyFile, data, 0644)
}
