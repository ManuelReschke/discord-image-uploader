package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"discord-image-uploader/internal/config"
	"discord-image-uploader/internal/discord"
	"discord-image-uploader/internal/history"
	"discord-image-uploader/internal/uploader"
	"discord-image-uploader/internal/watcher"
)

var (
	Version   = "dev"
	BuildTime = "unknown"
	GitCommit = "unknown"
)

func main() {
	configPath := flag.String("config", "config/config.json", "Path to configuration file")
	version := flag.Bool("version", false, "Show version information")
	flag.Parse()

	if *version {
		fmt.Printf("Discord Image Uploader\n")
		fmt.Printf("Version:    %s\n", Version)
		fmt.Printf("Built:      %s\n", BuildTime)
		fmt.Printf("Git Commit: %s\n", GitCommit)
		os.Exit(0)
	}

	log.Println("Starting Discord Image Uploader...")

	var err error
	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	var discordClient *discord.Client

	if cfg.Discord.WebhookURL != "" {
		discordClient, err = discord.NewWebhookClient(cfg.Discord.WebhookURL)
	} else {
		discordClient, err = discord.NewClient(cfg.Discord.Token, cfg.Discord.ChannelID)
	}

	if err != nil {
		log.Fatalf("Failed to create Discord client: %v", err)
	}
	defer discordClient.Close()

	err = discordClient.TestConnection(cfg.Discord.TestMessage, cfg.Discord.SendTestMessage)
	if err != nil {
		log.Fatalf("Failed to connect to Discord: %v", err)
	}

	uploadHistory, err := history.New(cfg.History.FilePath)
	if err != nil {
		log.Fatalf("Failed to create upload history: %v", err)
	}

	fileWatcher, err := watcher.New(
		cfg.Watcher.FolderPath,
		cfg.Watcher.SupportedFormats,
		cfg.Watcher.DeleteAfterUpload,
	)
	if err != nil {
		log.Fatalf("Failed to create file watcher: %v", err)
	}
	defer fileWatcher.Stop()

	imageUploader := uploader.New(cfg, discordClient, fileWatcher, uploadHistory)

	fileWatcher.Start()

	err = imageUploader.Start()
	if err != nil {
		log.Fatalf("Failed to start uploader: %v", err)
	}
	defer imageUploader.Stop()

	log.Println("Discord Image Uploader is running. Press Ctrl+C to stop.")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan

	log.Println("Shutting down gracefully...")
}
