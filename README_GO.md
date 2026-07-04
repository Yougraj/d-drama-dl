---

# Deha Drama Streamer (Go) `v0.0.1`

A high-performance Go port of the Terminal User Interface (TUI) and web application for searching, selecting, and streaming dramas. Built with Go, featuring both TUI and web interfaces.

![Version](https://img.shields.io/badge/version-0.0.1-blue)
![Go](https://img.shields.io/badge/go-1.21+-green)
![Platform](https://img.shields.io/badge/platform-Linux%20%7C%20macOS%20%7C%20Windows-orange)

## ✨ Features

- **Fuzzy Search:** Fast drama and episode searching using go-fuzzyfinder.
- **Clean TUI:** Cross-platform terminal UI with consistent header formatting.
- **Direct Streaming:** Extracts direct MP4/M3U8 links without ads or browser.
- **Subtitle Support:** Automatically detects and loads external subtitles.
- **Playback Controller:** Navigate between episodes with Next/Previous controls.
- **Web Interface:** Optional Flask-like web server with session management.
- **Cross-Platform:**
  - **Linux:** MPV integration.
  - **macOS:** Command-line video players.
  - **Windows:** Windows Media Player or external players.
  - **Android (Termux):** Native Intent triggers for MPV.
  - **iOS:** URL scheme triggers for VLC.

## 🛠 Prerequisites

### System Requirements
- **Go 1.21+** - Install from [golang.org](https://golang.org/dl)
- **Video Player:**
  - *Linux/Android:* [MPV](https://mpv.io/)
  - *macOS:* Built-in or MPV
  - *Windows:* External player or built-in
  - *iOS:* [VLC](https://www.videolan.org/vlc/download-ios.html)

### Optional: FZF for enhanced TUI
- Linux: `sudo apt install fzf`
- macOS: `brew install fzf`
- Termux: `pkg install fzf`

## 🚀 Installation & Setup

### Clone the Repository
```bash
git clone https://github.com/Yougraj/d-drama-dl.git
cd d-drama-dl
git checkout golang-conversion
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

### Run the Web Server
Edit `main.go` and uncomment the web server initialization:
```bash
go run app.go main.go
```

Then navigate to `http://localhost:5000` in your browser.

### Navigation Guide (TUI)
1. **Select Environment:** Choose your platform (Android, iOS, or Linux).
2. **Search:** Enter the drama name.
3. **Select Title:** Use arrow keys or type to filter results.
4. **Browse Episodes:** Pick the episode to start with.
5. **Playback Control:**
   - While the video plays, use the terminal control menu.
   - Select `NEXT >>` to immediately play the next episode.
   - Select `BACK TO EPISODES` to change the drama list.

## 📁 Project Structure

```text
.
├── main.go              # TUI Application Logic
├── app.go               # Web Server Implementation
├── go.mod               # Go module definition
├── go.sum               # Go dependencies checksum
├── README_GO.md         # Go-specific documentation
└── README.md            # Original Python documentation
```

## 🔄 Comparison: Python vs Go

| Feature | Python | Go |
|---------|--------|----|
| **Startup Time** | ~500ms | ~5ms |
| **Memory Usage** | 50-100MB | 5-10MB |
| **Binary Size** | N/A (interpreted) | ~15MB |
| **Concurrency** | Threading | Native goroutines |
| **Dependencies** | requests, beautifulsoup4, flask | goquery |
| **Cross-Compilation** | Limited | Native |

## 🔐 Security Notes

⚠️ **This is a development version** - Default credentials are hardcoded for testing. For production:
1. Use environment variables for credentials
2. Implement proper session management
3. Use HTTPS/TLS
4. Add CSRF protection
5. Sanitize all user inputs

## 📝 Building Executables

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

### Cross-Compilation
```bash
# Build for Linux from macOS/Windows
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o deha-streamer-linux

# Build for Windows from Linux/macOS
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o deha-streamer.exe
```

## 🧪 Testing

Run tests with:
```bash
go test -v ./...
```

## 📦 Dependencies

- **goquery** - HTML parsing (similar to BeautifulSoup)
- **go-fuzzyfinder** - Fuzzy search UI (similar to fzf)

## ⚠️ Disclaimer

This tool is for educational purposes only. It scrapes publicly available content from third-party providers. The developers do not host any content and are not responsible for the content retrieved through this tool. Users are solely responsible for ensuring they have the legal right to access and stream any content obtained through this application.

---

**Version 0.0.1 (Go)** | Converted from Python | Original by **Deha**
