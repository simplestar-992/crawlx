package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

var (
	maxDepth  = flag.Int("depth", 3, "Max crawl depth")
	maxPages  = flag.Int("max", 100, "Max pages to crawl")
	parallel  = flag.Int("p", 5, "Parallel workers")
	verbose   = flag.Bool("v", false, "Verbose")
	onlySame  = flag.Bool("same", true, "Only crawl same domain")
	userAgent = flag.String("ua", "CrawlX/1.0", "User agent")
)

type Page struct {
	URL   string
	Title string
	Links []string
}

func main() {
	flag.Parse()
	startURL := flag.Arg(0)

	if startURL == "" {
		fmt.Println("CrawlX - Web Crawler")
		fmt.Println("")
		fmt.Println("Usage:")
		fmt.Println("  crawlx https://example.com")
		fmt.Println("")
		flag.PrintDefaults()
		return
	}

	fmt.Printf("🕷️  CrawlX - Web Crawler\n")
	fmt.Printf("   Starting: %s\n", startURL)
	fmt.Printf("   Max pages: %d, Depth: %d, Parallel: %d\n", *maxPages, *maxDepth, *parallel)
	fmt.Println("")

	baseHost := getHost(startURL)
	visited := make(map[string]bool)
	var mu sync.Mutex
	pages := make(chan string, *maxPages)
	results := make(chan Page, *maxPages)

	var wg sync.WaitGroup

	for i := 0; i < *parallel; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for url := range pages {
				processPage(url, baseHost, visited, &mu, pages, results)
			}
		}()
	}

	pages <- startURL
	count := 0

	for result := range results {
		count++
		title := result.Title
		if title == "" {
			title = "(no title)"
		}
		if len(title) > 60 {
			title = title[:60] + "..."
		}
		fmt.Printf("📄 [%d] %s\n   %s\n", count, title, result.URL)

		if count >= *maxPages {
			break
		}
	}

	close(pages)
	wg.Wait()
	close(results)

	fmt.Printf("\n✅ Crawled %d pages\n", count)
}

func processPage(pageURL, baseHost string, visited map[string]bool, mu *sync.Mutex, pages chan string, results chan Page) {
	mu.Lock()
	if visited[pageURL] {
		mu.Unlock()
		return
	}
	visited[pageURL] = true
	mu.Unlock()

	client := &http.Client{Timeout: 10 * time.Second}
	req, _ := http.NewRequest("GET", pageURL, nil)
	req.Header.Set("User-Agent", *userAgent)

	resp, err := client.Do(req)
	if err != nil {
		if *verbose {
			fmt.Printf("❌ %s: %v\n", pageURL, err)
		}
		return
	}
	defer resp.Body.Close()

	page := Page{URL: pageURL}
	results <- page

	buf := make([]byte, 4096)
	resp.Body.Read(buf)
	body := string(buf)

	if idx := strings.Index(body, "<title>"); idx >= 0 {
		end := strings.Index(body[idx:], "</title>")
		if end > 0 {
			page.Title = strings.TrimSpace(body[idx+7 : idx+end])
		}
	}

	for _, link := range extractLinks(body, pageURL) {
		if getHost(link) != baseHost && *onlySame {
			continue
		}
		mu.Lock()
		if !visited[link] {
			select {
			case pages <- link:
			default:
			}
		}
		mu.Unlock()
	}
}

func extractLinks(html, baseURL string) []string {
	var links []string
	start := 0
	for {
		idx := strings.Index(html[start:], `href="`)
		if idx < 0 {
			break
		}
		idx += 6
		end := strings.Index(html[start+idx:], `"`)
		if end < 0 {
			break
		}
		link := html[start+idx : start+idx+end]
		if strings.HasPrefix(link, "http") {
			links = append(links, link)
		} else if strings.HasPrefix(link, "/") {
			base, _ := url.Parse(baseURL)
			links = append(links, base.Scheme+"://"+base.Host+link)
		}
		start += idx + end + 1
	}
	return links
}

func getHost(u string) string {
	if parsed, err := url.Parse(u); err == nil {
		return parsed.Host
	}
	return ""
}
