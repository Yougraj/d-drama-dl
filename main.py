import os
import re
import subprocess
import sys
from urllib.parse import quote

import requests
from bs4 import BeautifulSoup

# --- Configuration & Metadata ---
VERSION = "0.0.1"
BASE_URL = "https://kisskh.buzz/"
AJAX_URL = BASE_URL + "wp-admin/admin-ajax.php"
BLOGGER_BLOG_ID = "1422331367239821646"
BLOGGER_FEED_URL = f"https://www.blogger.com/feeds/{BLOGGER_BLOG_ID}/posts/default"


def clear():
    """Clears the terminal screen."""
    os.system("clear" if os.name != "nt" else "cls")


def draw_header(context=""):
    """Draws a consistent TUI header with versioning."""
    clear()
    # Header box with version number
    print("┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓")
    print(f"┃ DEHA DRAMA STREAMER (TUI)                     v{VERSION} ┃")
    if context:
        # Centers the context string within the 59-character width
        print(f"┃ {context.center(59)} ┃")
    print("┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛")


def call_fzf(options, prompt="Select"):
    """Helper to use fzf for selection. Returns the selected string."""
    try:
        input_str = "\n".join(options)
        process = subprocess.Popen(
            [
                "fzf",
                "--prompt",
                f"{prompt} > ",
                "--height",
                "40%",
                "--reverse",
                "--border",
                "--inline-info",
            ],
            stdin=subprocess.PIPE,
            stdout=subprocess.PIPE,
            stderr=None,
            text=True,
        )
        stdout, _ = process.communicate(input=input_str)
        return stdout.strip()
    except:
        return None


def play_video(url, sub_url, platform):
    """Triggers external players with subtitle support."""
    if not url:
        return
    try:
        if platform == "1":  # Android (MPV)
            cmd = [
                "am",
                "start",
                "--user",
                "0",
                "-a",
                "android.intent.action.VIEW",
                "-d",
                url,
                "-n",
                "io.mpv/.MPVActivity",
            ]
            if sub_url:
                cmd.extend(["--es", "subs", sub_url])
            subprocess.run(cmd, stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL)

        elif platform == "2":  # iOS (VLC)
            os.system(f"open vlc://{url}")

        elif platform == "3":  # Linux (MPV)
            cmd = ["mpv", url]
            if sub_url:
                cmd.append(f"--sub-file={sub_url}")
            # Run as a background process
            subprocess.Popen(cmd, stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL)
    except:
        pass


def get_search_results(query):
    """Search for dramas via AJAX."""
    payload = {
        "action": "fetch_live_movies",
        "keyword": query,
        "filter": "all",
        "page": "1",
        "is_popular": "0",
    }
    try:
        res = requests.post(AJAX_URL, data=payload, timeout=10)
        soup = BeautifulSoup(res.text, "html.parser")
        results = []
        for card in soup.select("a.movie-card"):
            title = card.select_one(".movie-title").get_text(strip=True)
            link = card["href"]
            ep = (
                card.select_one(".episode").get_text(strip=True)
                if card.select_one(".episode")
                else "Movie"
            )
            results.append({"title": title, "link": link, "display": f"{title} [{ep}]"})
        return results
    except:
        return []


def fetch_links(drama_title):
    """Extracts Video and Subtitle links from Blogger Feed."""
    feed_url = f"{BLOGGER_FEED_URL}?q={quote(drama_title)}&alt=json&max-results=1"
    try:
        data = requests.get(feed_url, timeout=10).json()
        if "entry" not in data["feed"]:
            return []
        content = data["feed"]["entry"][0]["content"]["$t"]
        eps = []
        for i, part in enumerate(content.split(";"), 1):
            if "|" in part:
                fields = part.split("|")
                v_url = fields[0].strip()
                s_url = ""
                if len(fields) > 2:
                    subs = fields[2].strip().split(",")
                    s_url = subs[0] if subs else ""
                if v_url.startswith("http"):
                    eps.append({"label": f"Episode {i}", "url": v_url, "sub": s_url})
        return eps
    except:
        return []


def playback_controller(episodes, start_idx, drama_title, platform):
    """TUI Loop for controlling playback (Next/Prev)."""
    current_idx = start_idx
    while True:
        ep = episodes[current_idx]
        sub_info = "(Subtitles Loaded)" if ep["sub"] else "(No Subtitles)"
        draw_header(f"Playing: {drama_title} > {ep['label']} {sub_info}")

        play_video(ep["url"], ep["sub"], platform)

        nav = []
        if current_idx < len(episodes) - 1:
            nav.append("NEXT >>")
        if current_idx > 0:
            nav.append("<< PREVIOUS")
        nav.extend(["REPLAY CURRENT", "BACK TO EPISODES", "EXIT TO SEARCH"])

        choice = call_fzf(nav, "Control")
        if not choice or "BACK TO" in choice:
            break
        if "EXIT TO" in choice:
            return "EXIT_TO_SEARCH"
        if "NEXT" in choice:
            current_idx += 1
        elif "PREVIOUS" in choice:
            current_idx -= 1
    return "CONTINUE"


def main():
    clear()
    print("==========================================")
    print(f"       DEHA DRAMA STREAMER v{VERSION}      ")
    print("==========================================")
    print("1. Android (MPV) | 2. iOS (VLC) | 3. Linux (MPV) | 4. URL Only")
    platform = input("Select Environment: ").strip()

    while True:
        draw_header("Main Search")
        query = input("Search Drama (or 'q' to quit): ").strip()
        if not query or query.lower() == "q":
            break

        results = get_search_results(query)
        if not results:
            continue

        drama_list = [r["display"] for r in results]
        draw_header(f"Results for: {query}")
        selected_display = call_fzf(drama_list, "Select Title")
        if not selected_display:
            continue

        selected_drama = results[drama_list.index(selected_display)]
        episodes = fetch_links(selected_drama["title"])
        if not episodes:
            continue

        while True:
            ep_list = [e["label"] for e in episodes]
            draw_header(f"Browse: {selected_drama['title']}")
            selected_ep = call_fzf(ep_list, "Select Episode")
            if not selected_ep:
                break

            idx = ep_list.index(selected_ep)
            status = playback_controller(
                episodes, idx, selected_drama["title"], platform
            )
            if status == "EXIT_TO_SEARCH":
                break

    clear()
    print(f"Deha Streamer v{VERSION} - Goodbye!")


if __name__ == "__main__":
    try:
        main()
    except KeyboardInterrupt:
        clear()
        sys.exit(0)
        sys.exit(0)
