# Discord Image Uploader

Ein Go-Tool, das automatisch Bilder aus einem bestimmten Ordner auf dem PC in einen Discord-Kanal hochlädt.

[![Latest Release](https://img.shields.io/github/v/release/ManuelReschke/discord-image-uploader?style=for-the-badge&logo=github)](https://github.com/ManuelReschke/discord-image-uploader/releases/latest)
[![Go Version](https://img.shields.io/github/go-mod/go-version/ManuelReschke/discord-image-uploader?style=for-the-badge&logo=go)](https://golang.org/)
[![License](https://img.shields.io/github/license/ManuelReschke/discord-image-uploader?style=for-the-badge)](LICENSE)
[![Downloads](https://img.shields.io/github/downloads/ManuelReschke/discord-image-uploader/total?style=for-the-badge&logo=github)](https://github.com/ManuelReschke/discord-image-uploader/releases)

## 📥 Download

Laden Sie die neueste Version für Ihr Betriebssystem herunter:

### Windows
- **[Windows 64-bit](https://github.com/ManuelReschke/discord-image-uploader/releases/latest/download/discord-image-uploader-windows-amd64.exe)** (Empfohlen)
- **[Windows 32-bit](https://github.com/ManuelReschke/discord-image-uploader/releases/latest/download/discord-image-uploader-windows-386.exe)**

### Linux
- **[Linux 64-bit (Intel/AMD)](https://github.com/ManuelReschke/discord-image-uploader/releases/latest/download/discord-image-uploader-linux-amd64)**
- **[Linux ARM64](https://github.com/ManuelReschke/discord-image-uploader/releases/latest/download/discord-image-uploader-linux-arm64)** (Raspberry Pi 4+)
- **[Linux 32-bit](https://github.com/ManuelReschke/discord-image-uploader/releases/latest/download/discord-image-uploader-linux-386)**

### macOS
- **[macOS Intel](https://github.com/ManuelReschke/discord-image-uploader/releases/latest/download/discord-image-uploader-darwin-amd64)**
- **[macOS Apple Silicon (M1/M2)](https://github.com/ManuelReschke/discord-image-uploader/releases/latest/download/discord-image-uploader-darwin-arm64)** (Empfohlen für neue Macs)

> **💡 Tipp:** Nach dem Download die Datei ausführbar machen:
> ```bash
> chmod +x discord-image-uploader-*
> ```

### 🚀 Schnellstart

1. **Binary herunterladen** (siehe oben)
2. **Konfiguration erstellen:**
   ```bash
   curl -O https://raw.githubusercontent.com/ManuelReschke/discord-image-uploader/main/config/config.example.json
   mv config.example.json config.json
   ```
3. **Discord Webhook erstellen** und URL in `config.json` eintragen
4. **Starten:**
   ```bash
   ./discord-image-uploader-* -config config.json
   ```

---

## Features

- ✅ **Automatische Ordnerüberwachung**: Überwacht einen konfigurierbaren Ordner auf neue Bilddateien in Echtzeit
- ✅ **Discord-Integration**: Automatisches Hochladen in einen Discord-Kanal über Bot-API oder Webhooks
- ✅ **Unterstützte Formate**: PNG, JPG, JPEG, GIF, WEBP
- ✅ **Batch-Upload**: Mehrere Bilder gleichzeitig hochladen
- ✅ **Konfigurierbar**: Upload-Intervalle, Batch-Größe, Dateigröße-Limits
- ✅ **Optional**: Löschen von Bildern nach erfolgreichem Upload
- ✅ **Robuste Fehlerbehandlung**: Comprehensive Logging und Graceful Shutdown
- ✅ **Dateigröße-Validierung**: Discord-konforme Größenlimits (8MB Standard)

## Installation

### Voraussetzungen

- Go 1.19 oder höher
- **Option 1**: Discord Bot Token und entsprechende Berechtigungen
- **Option 2**: Discord Webhook URL (einfacher Setup)

### Setup

1. **Repository klonen:**
   ```bash
   git clone https://github.com/ManuelReschke/discord-image-uploader.git
   cd discord-image-uploader
   ```

2. **Dependencies installieren:**
   ```bash
   go mod tidy
   ```

3. **Konfiguration erstellen:**
   ```bash
   cp config/config.example.json config/config.json
   ```

4. **Discord konfigurieren (wähle eine Option):**

   **Option A: Webhook (Empfohlen - Einfacher Setup)**
   - Gehe zu deinem Discord-Kanal
   - Rechtsklick → "Kanal bearbeiten" → "Integrationen" → "Webhooks"
   - Klicke "Neuer Webhook" und kopiere die Webhook-URL
   - Trage die Webhook-URL in `config/config.json` ein

   **Option B: Discord Bot**
   - Erstelle einen Discord Bot auf https://discord.com/developers/applications
   - Kopiere den Bot Token
   - Lade den Bot zu deinem Server ein mit "Send Messages" und "Attach Files" Berechtigungen
   - Finde die Channel ID des Zielkanals
   - Trage Token und Channel ID in `config/config.json` ein

5. **Tool kompilieren:**
   ```bash
   go build -o discord-image-uploader.exe ./cmd
   ```

## Konfiguration

Die Konfiguration erfolgt über eine JSON-Datei (`config/config.json`):

**Webhook-Konfiguration (Empfohlen):**
```json
{
  "discord": {
    "webhook_url": "https://discord.com/api/webhooks/YOUR_WEBHOOK_ID/YOUR_WEBHOOK_TOKEN"
  },
  "watcher": {
    "folder_path": "C:\\Users\\ManuelReschke\\Pictures\\Screenshots",
    "supported_formats": [".png", ".jpg", ".jpeg", ".gif", ".webp"],
    "delete_after_upload": false
  },
  "upload": {
    "batch_size": 5,
    "interval_seconds": 10,
    "max_file_size_mb": 8
  }
}
```

**Bot-Konfiguration (Alternative):**
```json
{
  "discord": {
    "token": "YOUR_BOT_TOKEN_HERE",
    "channel_id": "YOUR_CHANNEL_ID_HERE"
  },
  "watcher": {
    "folder_path": "C:\\Users\\ManuelReschke\\Pictures\\Screenshots",
    "supported_formats": [".png", ".jpg", ".jpeg", ".gif", ".webp"],
    "delete_after_upload": false
  },
  "upload": {
    "batch_size": 5,
    "interval_seconds": 10,
    "max_file_size_mb": 8
  }
}
```

### Konfigurationsoptionen

#### Discord-Konfiguration

| Parameter | Beschreibung | Erforderlich | Standard |
|-----------|--------------|--------------|----------|
| `discord.webhook_url` | Discord Webhook URL (Option A) | Webhook oder Bot | - |
| `discord.token` | Discord Bot Token (Option B) | Webhook oder Bot | - |
| `discord.channel_id` | Discord Channel ID (nur bei Bot) | Bei Bot-Token | - |

#### Weitere Optionen
| `watcher.folder_path` | Zu überwachender Ordner | - |
| `watcher.supported_formats` | Unterstützte Dateiformate | `[".png", ".jpg", ".jpeg", ".gif", ".webp"]` |
| `watcher.delete_after_upload` | Dateien nach Upload löschen | `false` |
| `upload.batch_size` | Anzahl Dateien pro Batch | `5` |
| `upload.interval_seconds` | Upload-Intervall in Sekunden | `10` |
| `upload.max_file_size_mb` | Maximale Dateigröße in MB | `8` |

## Verwendung

```bash
./discord-image-uploader.exe -config config/config.json
```

### Command Line Optionen

- `-config`: Pfad zur Konfigurationsdatei (Standard: `config/config.json`)

### Umgebungsvariablen

Konfigurationswerte können auch über Umgebungsvariablen gesetzt werden:

```bash
# Webhook-Konfiguration
export DISCORD_UPLOADER_DISCORD_WEBHOOK_URL="https://discord.com/api/webhooks/..."

# Bot-Konfiguration (Alternative)
export DISCORD_UPLOADER_DISCORD_TOKEN="your_token"
export DISCORD_UPLOADER_DISCORD_CHANNEL_ID="your_channel_id"

# Weitere Optionen
export DISCORD_UPLOADER_WATCHER_FOLDER_PATH="/path/to/watch"
```

## Projektstruktur

```
discord-image-uploader/
├── cmd/
│   └── main.go                 # Hauptanwendung
├── internal/
│   ├── config/
│   │   └── config.go          # Konfigurationsmanagement
│   ├── discord/
│   │   └── client.go          # Discord API Client
│   ├── watcher/
│   │   └── watcher.go         # File System Watcher
│   └── uploader/
│       └── uploader.go        # Upload-Logik
├── config/
│   └── config.example.json    # Beispielkonfiguration
├── go.mod
├── go.sum
└── README.md
```

## Funktionsweise

1. **Initialisierung**: Lädt Konfiguration und stellt Discord-Verbindung her (Bot oder Webhook)
2. **Ordnerüberwachung**: Überwacht den konfigurierten Ordner mit `fsnotify`
3. **Datei-Erkennung**: Erkennt neue Bilddateien in unterstützten Formaten
4. **Warteschlange**: Fügt Dateien einer Upload-Warteschlange hinzu
5. **Batch-Upload**: Lädt Dateien über Discord-API (Bot) oder HTTP-Requests (Webhook) hoch
6. **Cleanup**: Optional: Löscht Dateien nach erfolgreichem Upload

### Webhook vs. Bot

| Aspekt | Webhook | Bot |
|--------|---------|-----|
| **Setup** | Sehr einfach - nur URL kopieren | Komplex - Bot erstellen, Berechtigungen setzen |
| **Berechtigungen** | Automatisch verfügbar | Manuell konfigurieren |
| **Rate Limits** | Weniger restriktiv | Discord Bot Rate Limits |
| **Features** | Nur Datei-Upload | Erweiterte Discord-Features möglich |
| **Empfehlung** | ✅ Für einfache Uploads | Für erweiterte Bot-Funktionen |

## Entwicklung

### Dependencies

- [`github.com/bwmarrin/discordgo`](https://github.com/bwmarrin/discordgo) - Discord API Client
- [`github.com/fsnotify/fsnotify`](https://github.com/fsnotify/fsnotify) - File System Watcher
- [`github.com/spf13/viper`](https://github.com/spf13/viper) - Configuration Management

### Build

Verwenden Sie das Makefile für einfache Builds:

```bash
# Alle Plattformen
make build-all

# Nur aktuelle Plattform  
make build

# Development Build
make dev

# Vollständiger Release
make release

# Hilfe anzeigen
make help
```

**Manuelle Builds:**
```bash
# Development Build
go build -o discord-image-uploader ./cmd

# Cross-Platform Builds
make build-windows  # Windows builds
make build-linux    # Linux builds  
make build-mac      # macOS builds
```

### Tests (geplant)

```bash
go test ./...
```

## Troubleshooting

### Häufige Probleme

1. **"Failed to connect to Discord"**
   - **Bei Webhook**: Überprüfe die Webhook-URL auf Gültigkeit
   - **Bei Bot**: Überprüfe Bot Token und Channel ID, stelle sicher dass der Bot die nötigen Berechtigungen hat

2. **"Watch path does not exist"**
   - Überprüfe den Pfad in der Konfiguration
   - Stelle sicher, dass der Ordner existiert

3. **"File too large"**
   - Standard Discord-Limit ist 8MB
   - Nitro-Server haben 50MB Limit

### Logging

Das Tool protokolliert alle wichtigen Ereignisse:
- Verbindungsstatus
- Datei-Erkennung
- Upload-Status
- Fehler und Warnungen

## Lizenz

MIT License - siehe LICENSE Datei für Details.

## Contributing

1. Fork das Repository
2. Erstelle einen Feature Branch (`git checkout -b feature/amazing-feature`)
3. Committe die Änderungen (`git commit -m 'Add amazing feature'`)
4. Push zum Branch (`git push origin feature/amazing-feature`)
5. Öffne einen Pull Request

## Roadmap

- [ ] Bildkomprimierung vor Upload
- [ ] Web-Interface für Konfiguration
- [ ] Mehrere Discord-Server Support
- [ ] Statistiken und Monitoring
- [ ] Plugin-System für verschiedene Cloud-Provider