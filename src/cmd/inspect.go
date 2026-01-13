package cmd

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/spf13/cobra"
)

var inspectCmd = &cobra.Command{
	Use:   "inspect [url] [selector]",
	Short: "Inspect HTML elements at a URL using CSS selector",
	Long: `Inspect HTML elements at a given URL using a CSS selector.
This helps you understand the page structure and test selectors before scraping.

Examples:
  # Inspect h1 elements
  meilisearch-scraper inspect https://docs.example.com/page "h1"

  # Inspect navigation
  meilisearch-scraper inspect https://docs.example.com/page "nav.sidebar"

  # Inspect with class selector
  meilisearch-scraper inspect https://docs.example.com/page ".content h2"`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		inspectURL := args[0]
		selector := args[1]

		fetchURL := inspectURL
		if !strings.HasSuffix(inspectURL, ".html") {
			fetchURL = inspectURL + ".html"
		}

		log.Printf("Inspecting URL: %s (fetching: %s) with selector: %s", inspectURL, fetchURL, selector)

		resp, err := http.Get(fetchURL)
		if err != nil {
			log.Fatalf("Failed to fetch page: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Fatalf("Unexpected status code: %d", resp.StatusCode)
		}

		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			log.Fatalf("Failed to parse HTML: %v", err)
		}

		selection := doc.Find(selector)
		log.Printf("Found %d elements matching selector '%s'", selection.Length(), selector)
		fmt.Println()

		selection.Each(func(i int, s *goquery.Selection) {
			fmt.Printf("=== Element %d ===\n", i+1)
			fmt.Printf("Tag: %s\n", goquery.NodeName(s))

			if s.Length() > 0 {
				node := s.Get(0)
				if len(node.Attr) > 0 {
					fmt.Println("Attributes:")
					for _, attr := range node.Attr {
						fmt.Printf("  %s=\"%s\"\n", attr.Key, attr.Val)
					}
				}
			}

			text := strings.TrimSpace(s.Text())
			if text != "" {
				fmt.Printf("Text: %s\n", text)
			}

			html, _ := s.Html()
			if len(html) > 200 {
				html = html[:200] + "..."
			}
			fmt.Printf("HTML: %s\n", html)
			fmt.Println()
		})
	},
}
