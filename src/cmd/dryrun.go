package cmd

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/jansaidl/meilisearch-scraper/src"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var dryRunCmd = &cobra.Command{
	Use:   "dry-run [sitemap-url]",
	Short: "Scrape documentation and save to JSON file (no upload)",
	Long: `Scrape all URLs from a sitemap and save the extracted documents to data.json
instead of uploading to Meilisearch. Useful for testing and debugging.

Examples:
  # Dry run with sitemap URL
  meilisearch-scraper dry-run https://docs.example.com/sitemap.xml

  # Dry run with limit
  meilisearch-scraper dry-run https://docs.example.com/sitemap.xml --limit 5`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		sitemapURL := viper.GetString("sitemap.url")
		if len(args) > 0 {
			sitemapURL = args[0]
		}
		if sitemapURL == "" {
			log.Fatal("Sitemap URL is required (use argument or SITEMAP_URL env variable)")
		}

		limit, _ := cmd.Flags().GetInt("limit")

		configPath := viper.GetString("config")
		if configPath == "" {
			configPath = "config.json"
		}

		configFile, err := os.ReadFile(configPath)
		if err != nil {
			log.Fatalf("Failed to read config file %s: %v", configPath, err)
		}

		var config src.Config
		if err := json.Unmarshal(configFile, &config); err != nil {
			log.Fatalf("Failed to parse config file: %v", err)
		}

		log.Printf("Starting dry-run for sitemap: %s", sitemapURL)
		log.Println("Will save to data.json")

		sitemap, err := src.FetchSitemap(sitemapURL)
		if err != nil {
			log.Fatalf("Failed to fetch sitemap: %v", err)
		}

		log.Printf("Found %d URLs in sitemap", len(sitemap.URLs))

		urlsToProcess := sitemap.URLs
		if limit > 0 && limit < len(sitemap.URLs) {
			urlsToProcess = sitemap.URLs[:limit]
			log.Printf("Limiting to %d URLs", limit)
		}

		var documents []src.Document
		for i, url := range urlsToProcess {
			log.Printf("Scraping %d/%d: %s", i+1, len(urlsToProcess), url.Loc)

			docs, err := src.ScrapePage(url.Loc, &config)
			if err != nil {
				log.Printf("Failed to scrape %s: %v", url.Loc, err)
				continue
			}

			documents = append(documents, docs...)
			time.Sleep(200 * time.Millisecond)
		}

		log.Printf("Successfully scraped %d documents", len(documents))

		if len(documents) > 0 {
			data, err := json.MarshalIndent(documents, "", "  ")
			if err != nil {
				log.Fatalf("Failed to marshal documents: %v", err)
			}

			if err := os.WriteFile("data.json", data, 0644); err != nil {
				log.Fatalf("Failed to write data.json: %v", err)
			}

			log.Printf("Saved %d documents to data.json", len(documents))
		}
	},
}

func init() {
	dryRunCmd.Flags().Int("limit", 0, "Limit number of URLs to process (0 = no limit)")
}
