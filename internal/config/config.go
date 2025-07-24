package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Discord DiscordConfig `mapstructure:"discord"`
	Watcher WatcherConfig `mapstructure:"watcher"`
	Upload  UploadConfig  `mapstructure:"upload"`
	History HistoryConfig `mapstructure:"history"`
}

type DiscordConfig struct {
	WebhookURL      string `mapstructure:"webhook_url"`
	Token           string `mapstructure:"token"`
	ChannelID       string `mapstructure:"channel_id"`
	TestMessage     string `mapstructure:"test_message"`
	SendTestMessage bool   `mapstructure:"send_test_message"`
}

type WatcherConfig struct {
	FolderPath        string   `mapstructure:"folder_path"`
	SupportedFormats  []string `mapstructure:"supported_formats"`
	DeleteAfterUpload bool     `mapstructure:"delete_after_upload"`
}

type UploadConfig struct {
	BatchSize       int `mapstructure:"batch_size"`
	IntervalSeconds int `mapstructure:"interval_seconds"`
	MaxFileSizeMB   int `mapstructure:"max_file_size_mb"`
}

type HistoryConfig struct {
	FilePath            string `mapstructure:"file_path"`
	CleanupMissingFiles bool   `mapstructure:"cleanup_missing_files"`
}

func Load(configPath string) (*Config, error) {
	viper.SetConfigFile(configPath)
	viper.SetConfigType("json")

	viper.SetEnvPrefix("DISCORD_UPLOADER")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &config, nil
}

func validateConfig(config *Config) error {
	if config.Discord.WebhookURL == "" && config.Discord.Token == "" {
		return fmt.Errorf("either discord webhook URL or bot token is required")
	}

	if config.Discord.Token != "" && config.Discord.ChannelID == "" {
		return fmt.Errorf("discord channel ID is required when using bot token")
	}

	if config.Watcher.FolderPath == "" {
		return fmt.Errorf("watcher folder path is required")
	}

	if len(config.Watcher.SupportedFormats) == 0 {
		config.Watcher.SupportedFormats = []string{".png", ".jpg", ".jpeg", ".gif", ".webp"}
	}

	if config.Upload.BatchSize <= 0 {
		config.Upload.BatchSize = 5
	}

	if config.Upload.IntervalSeconds <= 0 {
		config.Upload.IntervalSeconds = 10
	}

	if config.Upload.MaxFileSizeMB <= 0 {
		config.Upload.MaxFileSizeMB = 8
	}

	if config.Discord.TestMessage == "" {
		config.Discord.TestMessage = "Test connection from Discord Image Uploader"
	}

	if config.History.FilePath == "" {
		config.History.FilePath = "data/upload_history.json"
	}

	return nil
}
