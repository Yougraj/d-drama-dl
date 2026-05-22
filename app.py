import json
import os
from urllib.parse import quote

import requests
from bs4 import BeautifulSoup
from flask import Flask, jsonify, redirect, render_template, request, session, url_for

app = Flask(__name__)
app.secret_key = "suku_love_secret_key_99"  # Change this to anything random

# Configuration
USER_CREDENTIALS = {"username": "Sudarshona", "password": "Suku#2005"}
HISTORY_FILE = "history.json"
BLOGGER_FEED_URL = "https://www.blogger.com/feeds/1422331367239821646/posts/default"
AJAX_URL = "https://kisskh.buzz/wp-admin/admin-ajax.php"

# --- Helper Functions ---


def load_history():
    if not os.path.exists(HISTORY_FILE):
        return []
    try:
        with open(HISTORY_FILE, "r") as f:
            return json.load(f)
    except:
        return []


def save_history(title, link, image):
    history = load_history()
    # Remove if already exists (to move to top)
    history = [item for item in history if item["link"] != link]
    history.insert(0, {"title": title, "link": link, "image": image})
    with open(HISTORY_FILE, "w") as f:
        json.dump(history[:20], f)  # Keep last 20


def scrape_search(query):
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
            results.append(
                {
                    "title": card.select_one(".movie-title").get_text(strip=True),
                    "link": card["href"],
                    "image": (
                        card.select_one("img")["src"] if card.select_one("img") else ""
                    ),
                    "ep": (
                        card.select_one(".episode").get_text(strip=True)
                        if card.select_one(".episode")
                        else "Movie"
                    ),
                }
            )
        return results
    except:
        return []


def fetch_episodes(drama_title):
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
                s_url = fields[2].strip().split(",")[0] if len(fields) > 2 else ""
                if v_url.startswith("http"):
                    eps.append({"label": f"EP {i}", "url": v_url, "sub": s_url})
        return eps
    except:
        return []


# --- Routes ---


@app.route("/login", methods=["GET", "POST"])
def login():
    if request.method == "POST":
        user = request.form.get("username")
        pw = request.form.get("password")
        if user == USER_CREDENTIALS["username"] and pw == USER_CREDENTIALS["password"]:
            session["logged_in"] = True
            session["user"] = user
            return redirect(url_for("index"))
        return render_template("login.html", error="Invalid Password, My Love!")
    return render_template("login.html")


@app.route("/")
def index():
    if not session.get("logged_in"):
        return redirect(url_for("login"))
    return render_template("index.html", history=load_history())


@app.route("/api/search")
def api_search():
    if not session.get("logged_in"):
        return jsonify([])
    query = request.args.get("q", "")
    return jsonify(scrape_search(query))


@app.route("/watch")
def watch():
    if not session.get("logged_in"):
        return redirect(url_for("login"))
    title = request.args.get("title")
    link = request.args.get("link")
    image = request.args.get("image")

    save_history(title, link, image)
    episodes = fetch_episodes(title)
    return render_template("player.html", title=title, episodes=episodes)


@app.route("/logout")
def logout():
    session.clear()
    return redirect(url_for("login"))


if __name__ == "__main__":
    # Use environment port for deployment (Render/Railway)
    port = int(os.environ.get("PORT", 5000))
    app.run(host="0.0.0.0", port=port)
