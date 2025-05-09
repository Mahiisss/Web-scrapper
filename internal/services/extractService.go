package services

import (
	"net/url"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

func ExtractService(page HTMLPage, mp *map[string]map[string]int) {
	invIndex := *mp
	parsedUrl, _ := url.Parse(page.URL)

	cleanText := ExtractTextFromHTML(page.HTML)

	words := extractWords(cleanText)

	for _, word := range words {
		if _, ok := invIndex[word]; !ok {
			invIndex[word] = make(map[string]int)
		}
		invIndex[word][parsedUrl.String()]++
	}
}

func ExtractTextFromHTML(htmlStr string) string {
	doc, err := html.Parse(strings.NewReader(htmlStr))
	if err != nil {
		return ""
	}

	var textBuilder strings.Builder
	var extract func(*html.Node)
	var insideBody bool // Flag to track if we're inside <body>

	extract = func(n *html.Node) {
		if n.Type == html.ElementNode {
			if n.Data == "body" {
				insideBody = true // Start extracting text
			} else if n.Data == "head" {
				insideBody = false // Stop extracting if inside <head>
			}
		}

		if insideBody && n.Type == html.TextNode {
			textBuilder.WriteString(n.Data)
			textBuilder.WriteString(" ") // Preserve spacing
		}

		// Recursively process child nodes
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			extract(c)
		}
	}

	extract(doc)
	return strings.TrimSpace(textBuilder.String())
}

func extractWords(text string) []string {
	words := strings.Fields(text)

	var result []string

	for _, word := range words {
		clean := CleanWord(word)
		if clean != "" {
			result = append(result, clean)
		}
	}

	return result
}

func CleanWord(word string) string {
	clean := strings.ToLower(word)
	re := regexp.MustCompile(`[^a-z0-9]`)
	return re.ReplaceAllString(clean, "")
}
