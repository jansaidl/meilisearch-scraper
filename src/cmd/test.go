package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/jansaidl/meilisearch-scraper/src"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var testCmd = &cobra.Command{
	Use:   "test [url]",
	Short: "Test scraping a single URL and print results",
	Long: `Test the scraping configuration on a single URL without uploading to Meilisearch.
This outputs the extracted documents in JSON format for inspection.

Examples:
  # Test single URL
  meilisearch-scraper test https://docs.example.com/getting-started

  # Test with custom config
  meilisearch-scraper test https://docs.example.com/api --config my-config.json`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		testURL := args[0]

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

		log.Printf("Testing scraping for URL: %s", testURL)
		docs, err := src.ScrapePage(testURL, &config)
		if err != nil {
			log.Fatalf("Failed to scrape page: %v", err)
		}

		data, err := json.MarshalIndent(docs, "", "  ")
		if err != nil {
			log.Fatalf("Failed to marshal documents: %v", err)
		}

		fmt.Println(string(data))
		log.Printf("Generated %d documents", len(docs))
	},
}
