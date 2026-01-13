package cmd

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/jansaidl/meilisearch-scraper/src"
	"github.com/meilisearch/meilisearch-go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all documents in Meilisearch index",
	Long: `List all documents stored in the Meilisearch index with their titles and URLs.
Displays hierarchy levels and document URLs for easy reference.

Examples:
  # List all documents from default index
  meilisearch-scraper list

  # List documents from specific index
  meilisearch-scraper list --index my-docs

  # Limit number of results
  meilisearch-scraper list --limit 20`,
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

		limit, _ := cmd.Flags().GetInt64("limit")
		if limit == 0 {
			limit = 100 // default limit
		}
		client := meilisearch.New(meilisearchURL, meilisearch.WithAPIKey(meilisearchKey))
		index := client.Index(indexName)

		// Get documents
		resp := &meilisearch.DocumentsResult{}
		err := index.GetDocuments(&meilisearch.DocumentsQuery{
			Limit: limit,
		}, resp)
		if err != nil {
			log.Fatalf("Failed to get documents: %v", err)
		}

		// Parse results
		var documents []src.Document
		jsonData, err := json.Marshal(resp.Results)
		if err != nil {
			log.Fatalf("Failed to marshal results: %v", err)
		}
		if err := json.Unmarshal(jsonData, &documents); err != nil {
			log.Fatalf("Failed to unmarshal documents: %v", err)
		}

		fmt.Printf("Index: %s\n", indexName)
		fmt.Printf("Total documents: %d\n", resp.Total)
		fmt.Printf("Showing: %d documents\n\n", len(documents))

		for i, doc := range documents {
			fmt.Printf("--- Document %d ---\n", i+1)
			fmt.Printf("ID: %s\n", doc.ObjectID)
			fmt.Printf("URL: %s\n", doc.URL)

			if doc.HierarchyLvl0 != nil {
				fmt.Printf("Lvl0: %s\n", *doc.HierarchyLvl0)
			}
			if doc.HierarchyLvl1 != nil {
				fmt.Printf("Lvl1: %s\n", *doc.HierarchyLvl1)
			}
			if doc.HierarchyLvl2 != nil {
				fmt.Printf("Lvl2: %s\n", *doc.HierarchyLvl2)
			}
			if doc.Anchor != "" {
				fmt.Printf("Anchor: %s\n", doc.Anchor)
			}
			if doc.Content != nil && len(*doc.Content) > 100 {
				fmt.Printf("Content: %s...\n", (*doc.Content)[:100])
			} else if doc.Content != nil {
				fmt.Printf("Content: %s\n", *doc.Content)
			}
			fmt.Println()
		}
	},
}

func init() {
	listCmd.Flags().Int64("limit", 100, "Maximum number of documents to list")
}
