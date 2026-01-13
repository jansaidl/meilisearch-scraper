package cmd

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/jansaidl/meilisearch-scraper/src"
	"github.com/meilisearch/meilisearch-go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var runCmd = &cobra.Command{
	Use:   "run [sitemap-url]",
	Short: "Scrape documentation and upload to Meilisearch",
	Long: `Scrape all URLs from a sitemap, extract content using configured CSS selectors,
and upload the documents to Meilisearch for full-text search.

The sitemap URL can be provided as argument or via SITEMAP_URL environment variable.

Examples:
  # Run with sitemap URL argument
  meilisearch-scraper run https://docs.example.com/sitemap.xml

  # Run with environment variable
  export SITEMAP_URL=https://docs.example.com/sitemap.xml
  meilisearch-scraper run

  # Limit number of URLs to process
  meilisearch-scraper run https://docs.example.com/sitemap.xml --limit 10`,
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

		meilisearchURL := viper.GetString("meilisearch.url")
		meilisearchKey := viper.GetString("meilisearch.key")
		indexName := viper.GetString("meilisearch.index")
		if indexName == "" {
			indexName = "docs"
		}

		if meilisearchURL == "" {
			log.Fatal("MEILISEARCH_HOST_URL is required")
		}
		if meilisearchKey == "" {
			log.Fatal("MEILISEARCH_API_KEY is required")
		}

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

		log.Printf("Starting scraper for sitemap: %s", sitemapURL)

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
			client := meilisearch.New(meilisearchURL, meilisearch.WithAPIKey(meilisearchKey))
			index := client.Index(indexName)

			log.Printf("Uploading documents to Meilisearch index: %s", indexName)
			task, err := index.AddDocuments(documents, nil)
			if err != nil {
				log.Fatalf("Failed to add documents: %v", err)
			}
			log.Printf("Upload task ID: %d", task.TaskUID)
		}

		log.Println("Scraping completed successfully")
	},
}

func init() {
	runCmd.Flags().Int("limit", 0, "Limit number of URLs to process (0 = no limit)")
}
