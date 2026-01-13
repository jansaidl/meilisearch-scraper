package cmd

import (
	"fmt"
	"log"

	"github.com/meilisearch/meilisearch-go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Show Meilisearch index statistics",
	Long: `Display statistics about the Meilisearch index including number of documents,
index size, and other metadata.

Examples:
  # Show stats for default index
  meilisearch-scraper stats

  # Show stats for specific index
  meilisearch-scraper stats --index my-docs`,
	Run: func(cmd *cobra.Command, args []string) {
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

		client := meilisearch.New(meilisearchURL, meilisearch.WithAPIKey(meilisearchKey))
		index := client.Index(indexName)

		// Get index stats
		stats, err := index.GetStats()
		if err != nil {
			log.Fatalf("Failed to get index stats: %v", err)
		}

		fmt.Printf("Index: %s\n", indexName)
		fmt.Printf("Number of documents: %d\n", stats.NumberOfDocuments)
		fmt.Printf("Is indexing: %v\n", stats.IsIndexing)
		fmt.Printf("\nField distribution:\n")
		for field, count := range stats.FieldDistribution {
			fmt.Printf("  %s: %d\n", field, count)
		}
	},
}
