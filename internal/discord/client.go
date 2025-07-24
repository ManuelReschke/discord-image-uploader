package discord

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/bwmarrin/discordgo"
)

type Client struct {
	session    *discordgo.Session
	channelID  string
	webhookURL string
}

func NewClient(token, channelID string) (*Client, error) {
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, fmt.Errorf("failed to create Discord session: %w", err)
	}

	err = session.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open Discord session: %w", err)
	}

	return &Client{
		session:   session,
		channelID: channelID,
	}, nil
}

func NewWebhookClient(webhookURL string) (*Client, error) {
	return &Client{
		webhookURL: webhookURL,
	}, nil
}

func (c *Client) Close() error {
	if c.session != nil {
		return c.session.Close()
	}
	return nil
}

func (c *Client) UploadImage(filePath string) error {
	if c.webhookURL != "" {
		return c.uploadImageViaWebhook(filePath)
	}
	return c.uploadImageViaBot(filePath)
}

func (c *Client) uploadImageViaBot(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", filePath, err)
	}
	defer file.Close()

	fileName := filepath.Base(filePath)

	_, err = c.session.ChannelFileSend(c.channelID, fileName, file)
	if err != nil {
		return fmt.Errorf("failed to upload file %s: %w", filePath, err)
	}

	log.Printf("Successfully uploaded: %s", fileName)
	return nil
}

func (c *Client) uploadImageViaWebhook(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", filePath, err)
	}
	defer file.Close()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	fileName := filepath.Base(filePath)
	part, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		return fmt.Errorf("failed to create form file: %w", err)
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return fmt.Errorf("failed to copy file data: %w", err)
	}

	err = writer.Close()
	if err != nil {
		return fmt.Errorf("failed to close multipart writer: %w", err)
	}

	req, err := http.NewRequest("POST", c.webhookURL, &body)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send webhook request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("webhook returned status %d", resp.StatusCode)
	}

	log.Printf("Successfully uploaded via webhook: %s", fileName)
	return nil
}

func (c *Client) UploadImages(filePaths []string) error {
	if c.webhookURL != "" {
		return c.uploadImagesViaWebhook(filePaths)
	}
	return c.uploadImagesViaBot(filePaths)
}

func (c *Client) uploadImagesViaBot(filePaths []string) error {
	var files []*discordgo.File

	for _, filePath := range filePaths {
		file, err := os.Open(filePath)
		if err != nil {
			log.Printf("Failed to open file %s: %v", filePath, err)
			continue
		}

		fileName := filepath.Base(filePath)
		files = append(files, &discordgo.File{
			Name:   fileName,
			Reader: file,
		})
	}

	if len(files) == 0 {
		return fmt.Errorf("no valid files to upload")
	}

	_, err := c.session.ChannelMessageSendComplex(c.channelID, &discordgo.MessageSend{
		Files: files,
	})

	for _, file := range files {
		if closer, ok := file.Reader.(interface{ Close() error }); ok {
			closer.Close()
		}
	}

	if err != nil {
		return fmt.Errorf("failed to upload batch: %w", err)
	}

	log.Printf("Successfully uploaded batch of %d files", len(files))
	return nil
}

func (c *Client) uploadImagesViaWebhook(filePaths []string) error {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	for _, filePath := range filePaths {
		file, err := os.Open(filePath)
		if err != nil {
			log.Printf("Failed to open file %s: %v", filePath, err)
			continue
		}
		defer file.Close()

		fileName := filepath.Base(filePath)
		part, err := writer.CreateFormFile("file", fileName)
		if err != nil {
			log.Printf("Failed to create form file for %s: %v", filePath, err)
			continue
		}

		_, err = io.Copy(part, file)
		if err != nil {
			log.Printf("Failed to copy file data for %s: %v", filePath, err)
			continue
		}
	}

	err := writer.Close()
	if err != nil {
		return fmt.Errorf("failed to close multipart writer: %w", err)
	}

	req, err := http.NewRequest("POST", c.webhookURL, &body)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send webhook request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("webhook returned status %d", resp.StatusCode)
	}

	log.Printf("Successfully uploaded batch of %d files via webhook", len(filePaths))
	return nil
}

func (c *Client) TestConnection(testMessage string, sendTest bool) error {
	if c.webhookURL != "" {
		return c.testWebhookConnection(testMessage, sendTest)
	}
	return c.testBotConnection()
}

func (c *Client) testBotConnection() error {
	_, err := c.session.Channel(c.channelID)
	if err != nil {
		return fmt.Errorf("failed to access channel %s: %w", c.channelID, err)
	}

	log.Printf("Successfully connected to Discord channel: %s", c.channelID)
	return nil
}

func (c *Client) testWebhookConnection(testMessage string, sendTest bool) error {
	if !sendTest {
		log.Printf("Skipping webhook test message")
		return nil
	}

	payload := map[string]string{
		"content": testMessage,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal test payload: %w", err)
	}

	req, err := http.NewRequest("POST", c.webhookURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create test request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send test webhook request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("webhook test returned status %d", resp.StatusCode)
	}

	log.Printf("Successfully tested webhook connection")
	return nil
}
