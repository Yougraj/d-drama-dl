package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/PuerkitoBio/goquery"
	fuzz "github.com/ktr0731/go-fuzzyfinder"
)

const (
	VERSION              = "0.0.1"
	BASE_URL             = "https://kisskh.buzz/"
	AJAX_URL             = BASE_URL + "wp-admin/admin-ajax.php"
	BLOGGER_BLOG_ID      = "1422331367239821646"
	BLOGGER_FEED_URL     = "https://www.blogger.com/feeds/" + BLOGGER_BLOG_ID + "/posts/default"
)

type Episode struct {
	Label string `json:"label"`
	URL   string `json:"url"`
	Sub   string `json:"sub"`
}

type SearchResult struct {
	Title   string `json:"title"`
	Link    string `json:"link"`
	Display string `json:"display"`
}

type BlogEntry struct {
	Content string `json:"$t"`
}

type BlogFeed struct {
	Feed struct {
		Entry []struct {
			Content BlogEntry `json:"content"`
		} `json:"entry"`
	} `json:"feed"`
}

var httpClient = &http.Client{
	Timeout: 10 * time.Second,
}

func clearScreen() {
	if runtime.GOOS == "windows" {
		exec.Command("cmd", "/c", "cls").Run()
	} else {
		exec.Command("clear").Run()
	}
}

func drawHeader(context string) {
	clearScreen()
	fmt.Println("┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓")
	fmt.Printf("┃ DEHA DRAMA STREAMER (TUI)                     v%s ┃\n", VERSION)
	if context != "" {
		padding := (59 - len(context)) / 2
		fmt.Printf("┃ %s%s%s ┃\n", strings.Repeat(" ", padding), context, strings.Repeat(" ", 59-len(context)-padding))
	}
	fmt.Println("┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛")
}

func getSearchResults(query string) []SearchResult {
	payload := url.Values{
		"action":      {"fetch_live_movies"},
		"keyword":     {query},
		"filter":      {"all"},
		"page":        {"1"},
		"is_popular":  {"0"},
	}

	resp, err := httpClient.PostForm(AJAX_URL, payload)
	if err != nil {
		return []SearchResult{}
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return []SearchResult{}
	}

	var results []SearchResult
	doc.Find("a.movie-card").Each(func(i int, s *goquery.Selection) {
		title := strings.TrimSpace(s.Find(".movie-title").Text())
		link, _ := s.Attr("href")
		ep := "Movie"
		if epText := strings.TrimSpace(s.Find(".episode").Text()); epText != "" {
			ep = epText
		}
		display := fmt.Sprintf("%s [%s]", title, ep)
		results = append(results, SearchResult{
			Title:   title,
			Link:    link,
			Display: display,
		})
	})

	return results
}

func fetchLinks(dramaTitle string) []Episode {
	feedURL := fmt.Sprintf("%s?q=%s&alt=json&max-results=1", BLOGGER_FEED_URL, url.QueryEscape(dramaTitle))

	resp, err := httpClient.Get(feedURL)
	if err != nil {
		return []Episode{}
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var blogFeed BlogFeed
	if err := json.Unmarshal(body, &blogFeed); err != nil {
		return []Episode{}
	}

	if len(blogFeed.Feed.Entry) == 0 {
		return []Episode{}
	}

	content := blogFeed.Feed.Entry[0].Content.Content
	var episodes []Episode

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
			episodes = append(episodes, Episode{
				Label: fmt.Sprintf("Episode %d", i+1),
				URL:   vURL,
				Sub:   sURL,
			})
		}
	}

	return episodes
}

func playVideo(videoURL, subURL, platform string) {
	if videoURL == "" {
		return
	}

	switch platform {
	case "1": // Android (MPV)
		cmd := exec.Command(
			"am", "start",
			"--user", "0",
			"-a", "android.intent.action.VIEW",
			"-d", videoURL,
			"-n", "io.mpv/.MPVActivity",
		)
		if subURL != "" {
			cmd.Args = append(cmd.Args, "--es", "subs", subURL)
		}
		cmd.Run()

	case "2": // iOS (VLC)
		exec.Command("open", fmt.Sprintf("vlc://%s", videoURL)).Run()

	case "3": // Linux (MPV)
		args := []string{videoURL}
		if subURL != "" {
			args = append(args, fmt.Sprintf("--sub-file=%s", subURL))
		}
		go exec.Command("mpv", args...).Run()
	}
}

func callFzf(options []string, prompt string) (string, error) {
	if len(options) == 0 {
		return "", fmt.Errorf("no options available")
	}

	idx, err := fuzz.Find(
		options,
		func(i int) string { return options[i] },
		fuzz.WithPreviewWindow(func(i, w, h int) string {
			if i < 0 {
				return ""
			}
			return options[i]
		}),
	)

	if err != nil {
		return "", err
	}

	return options[idx], nil
}

func playbackController(episodes []Episode, startIdx int, dramaTitle, platform string) string {
	currentIdx := startIdx

	for {
		ep := episodes[currentIdx]
		subInfo := "(No Subtitles)"
		if ep.Sub != "" {
			subInfo = "(Subtitles Loaded)"
		}
		drawHeader(fmt.Sprintf("Playing: %s > %s %s", dramaTitle, ep.Label, subInfo))

		playVideo(ep.URL, ep.Sub, platform)

		var nav []string
		if currentIdx < len(episodes)-1 {
			nav = append(nav, "NEXT >>")
		}
		if currentIdx > 0 {
			nav = append(nav, "<< PREVIOUS")
		}
		nav = append(nav, "REPLAY CURRENT", "BACK TO EPISODES", "EXIT TO SEARCH")

		choice, err := callFzf(nav, "Control")
		if err != nil || strings.Contains(choice, "BACK TO") {
			break
		}
		if strings.Contains(choice, "EXIT TO") {
			return "EXIT_TO_SEARCH"
		}
		if strings.Contains(choice, "NEXT") {
			currentIdx++
		} else if strings.Contains(choice, "PREVIOUS") {
			currentIdx--
		}
	}

	return "CONTINUE"
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			clearScreen()
			fmt.Printf("Deha Streamer v%s - Goodbye!\n", VERSION)
		}
	}()

	// Handle Ctrl+C gracefully
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		clearScreen()
		panic("interrupted")
	}()

	clearScreen()
	fmt.Println("==========================================")
	fmt.Printf("       DEHA DRAMA STREAMER v%s      \n", VERSION)
	fmt.Println("==========================================")
	fmt.Println("1. Android (MPV) | 2. iOS (VLC) | 3. Linux (MPV) | 4. URL Only")
	fmt.Print("Select Environment: ")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	platform := strings.TrimSpace(scanner.Text())

	for {
		drawHeader("Main Search")
		fmt.Print("Search Drama (or 'q' to quit): ")
		scanner.Scan()
		query := strings.TrimSpace(scanner.Text())

		if query == "" || query == "q" {
			break
		}

		results := getSearchResults(query)
		if len(results) == 0 {
			continue
		}

		dramaList := make([]string, len(results))
		for i, r := range results {
			dramaList[i] = r.Display
		}

		drawHeader(fmt.Sprintf("Results for: %s", query))
		selectedDisplay, err := callFzf(dramaList, "Select Title")
		if err != nil {
			continue
		}

		var selectedDrama SearchResult
		for _, r := range results {
			if r.Display == selectedDisplay {
				selectedDrama = r
				break
			}
		}

		episodes := fetchLinks(selectedDrama.Title)
		if len(episodes) == 0 {
			continue
		}

		for {
			epList := make([]string, len(episodes))
			for i, e := range episodes {
				epList[i] = e.Label
			}

			drawHeader(fmt.Sprintf("Browse: %s", selectedDrama.Title))
			selectedEp, err := callFzf(epList, "Select Episode")
			if err != nil {
				break
			}

			var idx int
			for i, e := range episodes {
				if e.Label == selectedEp {
					idx = i
					break
				}
			}

			status := playbackController(episodes, idx, selectedDrama.Title, platform)
			if status == "EXIT_TO_SEARCH" {
				break
			}
		}
	}

	clearScreen()
	fmt.Printf("Deha Streamer v%s - Goodbye!\n", VERSION)
}
