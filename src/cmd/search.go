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

var searchCmd = &cobra.Command{
	Use:   "search [query]",
	Short: "Search documents in Meilisearch index",
	Long: `Search for documents in the Meilisearch index using full-text search.
Returns matching documents with their titles, URLs, and content snippets.

Examples:
  # Search for documents containing "installation"
  meilisearch-scraper search "installation"

  # Search in specific index
  meilisearch-scraper search "api" --index my-docs

  # Limit number of results
  meilisearch-scraper search "configuration" --limit 5`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		query := args[0]

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
			limit = 10 // default limit for search
		}

		client := meilisearch.New(meilisearchURL, meilisearch.WithAPIKey(meilisearchKey))
		index := client.Index(indexName)

		// Perform search
		searchResp, err := index.Search(query, &meilisearch.SearchRequest{
			Limit: limit,
		})
		if err != nil {
			log.Fatalf("Failed to search documents: %v", err)
		}

		// Parse results
		var documents []src.Document
		jsonData, err := json.Marshal(searchResp.Hits)
		if err != nil {
			log.Fatalf("Failed to marshal results: %v", err)
		}
		if err := json.Unmarshal(jsonData, &documents); err != nil {
			log.Fatalf("Failed to unmarshal documents: %v", err)
		}

		fmt.Printf("Index: %s\n", indexName)
		fmt.Printf("Query: %s\n", query)
		fmt.Printf("Found: %d documents (showing %d)\n", searchResp.EstimatedTotalHits, len(documents))
		fmt.Printf("Processing time: %dms\n\n", searchResp.ProcessingTimeMs)

		if len(documents) == 0 {
			fmt.Println("No results found.")
			return
		}

		for i, doc := range documents {
			fmt.Printf("--- Result %d ---\n", i+1)
			fmt.Printf("ID: %s\n", doc.ObjectID)
			fmt.Printf("URL: %s\n", doc.URL)

			// Display hierarchy
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

			// Display content snippet
			if doc.Content != nil {
				content := *doc.Content
				if len(content) > 200 {
					fmt.Printf("Content: %s...\n", content[:200])
				} else {
					fmt.Printf("Content: %s\n", content)
				}
			}
			fmt.Println()
		}
	},
}

func init() {
	searchCmd.Flags().Int64("limit", 10, "Maximum number of results to return")
}
