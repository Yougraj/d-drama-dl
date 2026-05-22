---

# Deha Drama Streamer `v0.0.1`

A high-performance, clean Terminal User Interface (TUI) for searching, selecting, and streaming dramas directly from the terminal. Built with Python and powered by `fzf` for a lightning-fast fuzzy-finding experience.

![Version](https://img.shields.io/badge/version-0.0.1-blue)
![Python](https://img.shields.io/badge/python-3.10+-green)
![Platform](https://img.shields.io/badge/platform-Linux%20%7C%20Android%20%7C%20iOS-orange)

## ✨ Features

- **Fuzzy Search:** Instant drama and episode searching using `fzf`.
- **Clean TUI:** Persistent headers and auto-wiping terminal to ensure zero clutter.
- **Direct Streaming:** Extracts direct MP4/M3U8 links (no ads, no browser needed).
- **Subtitle Support:** Automatically detects and loads external subtitles.
- **Playback Controller:** Navigate between **Next** and **Previous** episodes directly from the TUI.
- **Cross-Platform:** Optimized for:
  - **Linux:** MPV integration.
  - **Android (Termux):** Native Intent triggers for MPV.
  - **iOS (a-Shell/Pyto):** URL scheme triggers for VLC.

## 🛠 Prerequisites

Before running the script, ensure you have the following installed:

### 1. System Requirements
- **fzf**: The fuzzy finder must be installed on your system.
  - *Linux:* `sudo apt install fzf`
  - *Android (Termux):* `pkg install fzf`
- **Video Player**:
  - *Linux/Android:* [MPV](https://mpv.io/)
  - *iOS:* [VLC](https://www.videolan.org/vlc/download-ios.html)

### 2. Python Environment
This project is optimized for [**uv**](https://github.com/astral-sh/uv).

## 🚀 Installation & Setup

```bash
# Clone the repository
git clone https://github.com/yourusername/deha-drama-streamer.git
cd deha-drama-streamer

# Sync dependencies using uv
uv sync
```

## 📖 Usage

Run the script using `uv`:

```bash
uv run main.py
```

### Navigation Guide:
1. **Select Environment:** Choose your platform (Linux, Android, or iOS).
2. **Search:** Enter the name of the drama.
3. **Select Title:** Use arrow keys or type to filter results in `fzf`.
4. **Browse Episodes:** Pick the episode you want to start with.
5. **Playback Control:**
   - While the video plays, the terminal remains in "Control Mode."
   - Select `NEXT >>` to immediately push the next episode to your player.
   - Select `BACK TO EPISODES` to change the current drama list.

## 📱 Mobile Setup

### Android (Termux)
1. Install Termux from F-Droid.
2. Install dependencies: `pkg install python fzf`.
3. Install MPV from the Play Store.
4. Run the script and choose option `1`.

### iOS (a-Shell)
1. Install **a-Shell** and **VLC** from the App Store.
2. Install requirements: `pip install requests beautifulsoup4`.
3. Run the script and choose option `2`.

## ⚙️ Project Structure

```text
.
├── main.py              # Main TUI Logic
├── pyproject.toml       # Project metadata and dependencies
├── uv.lock              # Locked dependency tree
└── README.md            # You are here!
```

## ⚠️ Disclaimer
This tool is for educational purposes only. It scrapes publicly available content from third-party providers. The developers do not host any content and are not responsible for the content retrieved.

---

**Version 0.0.1** | Developed by **Deha**
