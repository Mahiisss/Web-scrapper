package services

import (
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/html"
)

const maxDepth = 5 // Stop recursion after 5 levels

type CrawlService struct {
	baseURL  string
	visited  map[string]bool
	mu       sync.Mutex
	wg       sync.WaitGroup
	htmlChan chan HTMLPage // Channel to send crawled HTML to ExtractService
}

type HTMLPage struct {
	URL  string
	HTML string
}

func NewCrawlService(baseURL string, htmlChan chan HTMLPage) *CrawlService {
	return &CrawlService{
		baseURL:  baseURL,
		visited:  make(map[string]bool),
		htmlChan: htmlChan,
	}
}

func (c *CrawlService) Start() {
	c.wg.Add(1)
	go c.crawlLinks(c.baseURL, c.baseURL, 0)
	c.wg.Wait()
	close(c.htmlChan) // Close the channel when crawling is done
}

func (c *CrawlService) crawlLinks(baseURL string, url string, depth int) {
	defer c.wg.Done()

	c.mu.Lock()
	if c.visited[url] || depth > maxDepth {
		c.mu.Unlock()
		return
	}
	c.visited[url] = true
	c.mu.Unlock()

	body := extractBody(url)
	if body == "" {
		return
	}

	// Send the crawled HTML to the ExtractService via the channel
	c.htmlChan <- HTMLPage{URL: url, HTML: body}

	links := extractLinks(body, baseURL)
	for _, l := range links {
		c.wg.Add(1)
		go c.crawlLinks(baseURL, l, depth+1)
	}
}

func extractLinks(body string, baseURL string) []string {
	var links []string
	tokenizer := html.NewTokenizer(strings.NewReader(body))

	for {
		tokenType := tokenizer.Next()
		switch tokenType {
		case html.ErrorToken:
			return links
		case html.StartTagToken:
			token := tokenizer.Token()
			if token.Data == "a" {
				for _, attr := range token.Attr {
					if attr.Key == "href" {
						link := Clean(baseURL, attr.Val)
						if link != "" {
							links = append(links, link)
						}
					}
				}
			}
		}
	}
}

func extractBody(url string) string {
	client := &http.Client{Timeout: 10 * time.Second}

	resp, err := client.Get(url)
	if err != nil {
		log.Println("Error accessing URL: ", err)
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		//log.Printf("Failed to fetch page: %d - %s", resp.StatusCode, resp.Status)
		return ""
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response body: ", err)
		return ""
	}

	return string(body)
}
