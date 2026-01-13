package cmd

import (
	"log"

	"github.com/meilisearch/meilisearch-go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete all documents from Meilisearch index",
	Long: `Delete all documents from the specified Meilisearch index.
This is useful when you want to start fresh with new data.

Examples:
  # Delete all documents from default index
  meilisearch-scraper delete

  # Delete from specific index
  meilisearch-scraper delete --index my-docs`,
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

		log.Printf("Deleting all documents from index: %s", indexName)
		task, err := index.DeleteAllDocuments(nil)
		if err != nil {
			log.Fatalf("Failed to delete documents: %v", err)
		}
		log.Printf("Delete task ID: %d", task.TaskUID)
		log.Println("All documents deleted successfully")
	},
}
