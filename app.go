package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const (
	HISTORY_FILE_WEB = "history.json"
	BLOGGER_FEED_WEB = "https://www.blogger.com/feeds/1422331367239821646/posts/default"
	AJAX_URL_WEB     = "https://kisskh.buzz/wp-admin/admin-ajax.php"
)

type HistoryItem struct {
	Title string `json:"title"`
	Link  string `json:"link"`
	Image string `json:"image"`
}

type SearchItemWeb struct {
	Title string `json:"title"`
	Link  string `json:"link"`
	Image string `json:"image"`
	Ep    string `json:"ep"`
}

type EpisodeWeb struct {
	Label string `json:"label"`
	URL   string `json:"url"`
	Sub   string `json:"sub"`
}

type Server struct {
	SecretKey      string
	Username       string
	Password       string
	SessionStorage map[string]string // Simple in-memory session storage (username -> session_token)
}

func loadHistory() []HistoryItem {
	data, err := os.ReadFile(HISTORY_FILE_WEB)
	if err != nil {
		return []HistoryItem{}
	}

	var history []HistoryItem
	json.Unmarshal(data, &history)
	return history
}

func saveHistory(title, link, image string) error {
	history := loadHistory()

	// Remove if already exists (to move to top)
	var filtered []HistoryItem
	for _, item := range history {
		if item.Link != link {
			filtered = append(filtered, item)
		}
	}

	// Insert at beginning
	newHistory := []HistoryItem{
		{Title: title, Link: link, Image: image},
	}
	newHistory = append(newHistory, filtered...)

	// Keep last 20
	if len(newHistory) > 20 {
		newHistory = newHistory[:20]
	}

	data, _ := json.Marshal(newHistory)
	return os.WriteFile(HISTORY_FILE_WEB, data, 0644)
}

func scrapeSearch(query string) []SearchItemWeb {
	payload := url.Values{
		"action":      {"fetch_live_movies"},
		"keyword":     {query},
		"filter":      {"all"},
		"page":        {"1"},
		"is_popular":  {"0"},
	}

	resp, err := http.PostForm(AJAX_URL_WEB, payload)
	if err != nil {
		return []SearchItemWeb{}
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return []SearchItemWeb{}
	}

	var results []SearchItemWeb
	doc.Find("a.movie-card").Each(func(i int, s *goquery.Selection) {
		title := strings.TrimSpace(s.Find(".movie-title").Text())
		link, _ := s.Attr("href")
		image := ""
		if img := s.Find("img"); img.Length() > 0 {
			image, _ = img.Attr("src")
		}
		ep := "Movie"
		if epText := strings.TrimSpace(s.Find(".episode").Text()); epText != "" {
			ep = epText
		}
		results = append(results, SearchItemWeb{
			Title: title,
			Link:  link,
			Image: image,
			Ep:    ep,
		})
	})

	return results
}

func fetchEpisodes(dramaTitle string) []EpisodeWeb {
	feedURL := fmt.Sprintf("%s?q=%s&alt=json&max-results=1", BLOGGER_FEED_WEB, url.QueryEscape(dramaTitle))

	resp, err := http.Get(feedURL)
	if err != nil {
		return []EpisodeWeb{}
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var blogFeed BlogFeed
	if err := json.Unmarshal(body, &blogFeed); err != nil {
		return []EpisodeWeb{}
	}

	if len(blogFeed.Feed.Entry) == 0 {
		return []EpisodeWeb{}
	}

	content := blogFeed.Feed.Entry[0].Content.Content
	var episodes []EpisodeWeb

	parts := strings.Split(content, ";")
	for i, part := range parts {
		if !strings.Contains(part, "|") {
			continue
		}

		fields := strings.Split(part, "|")
		vURL := strings.TrimSpace(fields[0])
		sURL := ""

		if len(fields) > 2 {
			subFields := strings.Split(strings.TrimSpace(fields[2]), ",")
			if len(subFields) > 0 {
				sURL = subFields[0]
			}
		}

		if strings.HasPrefix(vURL, "http") {
			episodes = append(episodes, EpisodeWeb{
				Label: fmt.Sprintf("EP %d", i+1),
				URL:   vURL,
				Sub:   sURL,
			})
		}
	}

	return episodes
}

func (s *Server) generateSessionToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func (s *Server) getSessionFromCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie("session")
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

func (s *Server) setSessionCookie(w http.ResponseWriter, username, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLax,
	})
	s.SessionStorage[token] = username
}

func (s *Server) login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.Header().Set("Content-Type", "text/html")
		html := `<!DOCTYPE html>
<html>
<head><title>Login - Deha Drama Streamer</title></head>
<body>
<h2>Deha Drama Streamer Login</h2>
<form method="POST">
<input type="text" name="username" placeholder="Username" required>
<input type="password" name="password" placeholder="Password" required>
<button type="submit">Login</button>
</form>
</body>
</html>`
		w.Write([]byte(html))
		return
	}

	if r.Method == "POST" {
		r.ParseForm()
		username := r.FormValue("username")
		password := r.FormValue("password")

		if username == s.Username && password == s.Password {
			token := s.generateSessionToken()
			s.setSessionCookie(w, username, token)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		w.Header().Set("Content-Type", "text/html")
		html := `<!DOCTYPE html>
<html>
<head><title>Login - Deha Drama Streamer</title></head>
<body>
<h2>Deha Drama Streamer Login</h2>
<p style="color:red;">Invalid Password, My Love!</p>
<form method="POST">
<input type="text" name="username" placeholder="Username" required>
<input type="password" name="password" placeholder="Password" required>
<button type="submit">Login</button>
</form>
</body>
</html>`
		w.Write([]byte(html))
	}
}

func (s *Server) isLoggedIn(r *http.Request) bool {
	token, err := s.getSessionFromCookie(r)
	if err != nil {
		return false
	}
	_, exists := s.SessionStorage[token]
	return exists
}

func (s *Server) index(w http.ResponseWriter, r *http.Request) {
	if !s.isLoggedIn(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	history := loadHistory()
	w.Header().Set("Content-Type", "text/html")

	var historyHTML string
	for _, item := range history {
		historyHTML += fmt.Sprintf(`<div><img src="%s" width="100"><a href="/watch?title=%s&link=%s&image=%s">%s</a></div>\n`,
			item.Image, url.QueryEscape(item.Title), url.QueryEscape(item.Link), url.QueryEscape(item.Image), item.Title)
	}

	html := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head><title>Deha Drama Streamer</title></head>
<body>
<h1>Deha Drama Streamer</h1>
<input id="searchInput" type="text" placeholder="Search drama...">
<div id="results"></div>
<h2>Watch History</h2>
<div id="history">%s</div>
<a href="/logout">Logout</a>
<script>
document.getElementById('searchInput').addEventListener('keyup', function() {
  var query = this.value;
  fetch('/api/search?q=' + encodeURIComponent(query))
    .then(r => r.json())
    .then(data => {
      var html = '';
      data.forEach(item => {
        html += '<div><img src="' + item.image + '" width="100"><a href="/watch?title=' + encodeURIComponent(item.title) + '&link=' + encodeURIComponent(item.link) + '&image=' + encodeURIComponent(item.image) + '">\n' + item.title + ' [' + item.ep + ']</a></div>';
      });
      document.getElementById('results').innerHTML = html;
    });
});
</script>
</body>
</html>`, historyHTML)
	w.Write([]byte(html))
}

func (s *Server) apiSearch(w http.ResponseWriter, r *http.Request) {
	if !s.isLoggedIn(r) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("[]"))
		return
	}

	query := r.URL.Query().Get("q")
	results := scrapeSearch(query)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func (s *Server) watch(w http.ResponseWriter, r *http.Request) {
	if !s.isLoggedIn(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	title := r.URL.Query().Get("title")
	link := r.URL.Query().Get("link")
	image := r.URL.Query().Get("image")

	saveHistory(title, link, image)
	episodes := fetchEpisodes(title)

	w.Header().Set("Content-Type", "text/html")

	var epHTML string
	for _, ep := range episodes {
		epHTML += fmt.Sprintf(`<option value="%s">%s</option>\n`, ep.URL, ep.Label)
	}

	html := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head><title>%s - Deha Drama Streamer</title></head>
<body>
<h1>%s</h1>
<select id="episodeSelect">
%s
</select>
<button onclick="playEpisode()">Play</button>
<video id="videoPlayer" width="100%%" controls></video>
<script>
function playEpisode() {
  var url = document.getElementById('episodeSelect').value;
  document.getElementById('videoPlayer').src = url;
}
</script>
<a href="/">Back</a>
</body>
</html>`, title, title, epHTML)
	w.Write([]byte(html))
}

func (s *Server) logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	})
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func startWebServer(port int) {
	server := &Server{
		Username:       "Sudarshona",
		Password:       "Suku#2005",
		SessionStorage: make(map[string]string),
	}

	http.HandleFunc("/login", server.login)
	http.HandleFunc("/", server.index)
	http.HandleFunc("/api/search", server.apiSearch)
	http.HandleFunc("/watch", server.watch)
	http.HandleFunc("/logout", server.logout)

	addr := fmt.Sprintf("0.0.0.0:%d", port)
	fmt.Printf("Server running on http://0.0.0.0:%d\n", port)
	http.ListenAndServe(addr, nil)
}
