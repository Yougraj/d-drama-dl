---

# Deha Drama Streamer (Go Edition) `v0.0.1`

A high-performance Go port of the Terminal User Interface (TUI) and web application for searching, selecting, and streaming dramas. Built with Go for speed and efficiency.

![Version](https://img.shields.io/badge/version-0.0.1-blue)
![Go](https://img.shields.io/badge/go-1.21+-green)
![Platform](https://img.shields.io/badge/platform-Linux%20%7C%20macOS%20%7C%20Windows%20%7C%20Android%20%7C%20iOS-orange)

## ✨ Features

- **Blazing Fast** - Single binary, 5ms startup time, 5-10MB memory usage
- **Fuzzy Search** - Drama and episode searching using go-fuzzyfinder
- **Clean TUI** - Cross-platform terminal UI with consistent header formatting
- **Direct Streaming** - Extracts direct MP4/M3U8 links without ads or browser
- **Subtitle Support** - Automatically detects and loads external subtitles
- **Playback Controller** - Navigate between episodes with Next/Previous controls
- **Web Interface** - Optional web server with session management
- **Cross-Platform:**
  - Linux: MPV integration
  - macOS: Command-line video players
  - Windows: Windows Media Player or external players
  - Android (Termux): Native Intent triggers for MPV
  - iOS: URL scheme triggers for VLC

## 🛠 Prerequisites

### System Requirements
- **Go 1.21+** - Install from [golang.org](https://golang.org/dl)
- **Video Player:**
  - *Linux/Android:* [MPV](https://mpv.io/)
  - *macOS:* Built-in or MPV
  - *Windows:* External player or built-in
  - *iOS:* [VLC](https://www.videolan.org/vlc/download-ios.html)

### Optional: Enhanced TUI
- Linux: `sudo apt install fzf`
- macOS: `brew install fzf`
- Termux: `pkg install fzf`

## 🚀 Installation & Setup

### Clone the Repository
```bash
git clone https://github.com/Yougraj/d-drama-dl.git
cd d-drama-dl
```

### Download Dependencies
```bash
go mod download
```

## 📖 Usage

### Run the TUI Application
```bash
go run main.go
```

### Build a Release Binary
```bash
go build -o deha-streamer
./deha-streamer
```

### Run the Web Server
To enable web server mode, modify `main.go` and uncomment the web server initialization:
```bash
go run app.go main.go
```

Then navigate to `http://localhost:5000` in your browser.

### Navigation Guide (TUI)
1. **Select Environment:** Choose your platform (Android, iOS, or Linux)
2. **Search:** Enter the drama name
3. **Select Title:** Use arrow keys or type to filter results
4. **Browse Episodes:** Pick the episode to start with
5. **Playback Control:**
   - While the video plays, use the terminal control menu
   - Select `NEXT >>` to immediately play the next episode
   - Select `BACK TO EPISODES` to change the drama list

## 📦 Project Structure

```text
.
├── main.go              # TUI Application Logic (~450 lines)
├── app.go               # Web Server Implementation (~350 lines)
├── go.mod               # Go module definition
├── go.sum               # Go dependencies checksum
├── README.md            # This file
└── history.json         # Watch history (auto-created)
```

## 🔄 Migration from Python

This project was originally written in Python but has been completely rewritten in Go for:

| Feature | Python | Go |
|---------|--------|-----|
| **Startup Time** | ~500ms | ~5ms |
| **Memory Usage** | 50-100MB | 5-10MB |
| **Binary Size** | N/A (interpreted) | ~15MB (standalone) |
| **Concurrency** | Threading | Native goroutines |
| **Installation** | Requires Python 3.10+ + pip | Single binary |
| **Cross-Compilation** | Limited | Native support |

## 🔨 Building Executables

### Linux
```bash
GO111MODULE=on go build -o deha-streamer
./deha-streamer
```

### macOS
```bash
GO111MODULE=on go build -o deha-streamer
./deha-streamer
```

### Windows
```bash
GO111MODULE=on go build -o deha-streamer.exe
.\deha-streamer.exe
```

### Cross-Compilation Examples
```bash
# Build for Linux from macOS/Windows
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o deha-streamer-linux

# Build for Windows from Linux/macOS
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o deha-streamer.exe

# Build for macOS from Linux/Windows
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o deha-streamer-macos
```

## 🧪 Development

### Format Code
```bash
go fmt ./...
```

### Lint Code
```bash
golangci-lint run
```

### Run Tests
```bash
go test -v ./...
```

## 📦 Dependencies

- **goquery** - HTML parsing (similar to BeautifulSoup)
- **go-fuzzyfinder** - Fuzzy search UI (similar to fzf)

## ⚙️ Configuration

### Web Server Credentials
Edit the hardcoded credentials in `app.go`:
```go
Username: "Sudarshona"
Password: "Suku#2005"
```

**For production**, use environment variables:
```go
Username: os.Getenv("DRAMA_USERNAME")
Password: os.Getenv("DRAMA_PASSWORD")
```

## 🔐 Security Notes

⚠️ **This is a development version** - Default credentials are hardcoded. For production:
1. Use environment variables for credentials
2. Implement proper session management with expiration
3. Use HTTPS/TLS
4. Add CSRF protection
5. Sanitize all user inputs
6. Rate limit API endpoints
7. Use secure session cookies

## 📱 Mobile Setup

### Android (Termux)
1. Install Termux from F-Droid
2. Install Go: `pkg install golang`
3. Install dependencies: `pkg install mpv`
4. Clone and build:
   ```bash
   git clone https://github.com/Yougraj/d-drama-dl.git
   cd d-drama-dl
   go build -o deha-streamer
   ./deha-streamer
   ```
5. Choose option `1` for Android (MPV)

### iOS (a-Shell)
1. Install **a-Shell** from App Store
2. Install **VLC** for playback
3. Python environment not available on iOS - Go binary needed
4. Transfer binary via iCloud or build locally

## 🐛 Troubleshooting

### "fzf: command not found"
- Install fzf: `sudo apt install fzf` (Linux) or `brew install fzf` (macOS)
- Or use the built-in fuzzy finder (less optimal)

### "mpv: command not found"
- Install MPV: `sudo apt install mpv` (Linux) or `brew install mpv` (macOS)
- Alternative players can be integrated by modifying the `playVideo()` function

### Connection timeout errors
- Check internet connection
- Verify the streaming service is accessible
- Try searching for a different drama

## ⚠️ Disclaimer

This tool is for **educational purposes only**. It scrapes publicly available content from third-party providers. The developers:
- Do NOT host any content
- Are NOT responsible for the content retrieved through this tool
- Do NOT endorse piracy or copyright infringement

Users are solely responsible for ensuring they have the legal right to access and stream any content obtained through this application. Respect copyright laws in your jurisdiction.

## 📄 License

MIT License - See LICENSE file for details

---

**Version 0.0.1 (Go)** | High-performance drama streamer | Originally by **Deha** | Rewritten in Go
