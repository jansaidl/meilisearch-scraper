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

var detailCmd = &cobra.Command{
	Use:   "detail [document-id]",
	Short: "Show full document details by ID",
	Long: `Display complete information about a specific document in the Meilisearch index
using its objectID. Shows all fields including full content.

Examples:
  # Show document details by ID
  meilisearch-scraper detail abc123def456...

  # Show document from specific index
  meilisearch-scraper detail abc123def456... --index my-docs`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		documentID := args[0]

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

		// Get specific document
		var result map[string]interface{}
		err := index.GetDocument(documentID, nil, &result)
		if err != nil {
			log.Fatalf("Failed to get document: %v", err)
		}

		// Parse to Document struct for structured display
		jsonData, err := json.Marshal(result)
		if err != nil {
			log.Fatalf("Failed to marshal result: %v", err)
		}

		var doc src.Document
		if err := json.Unmarshal(jsonData, &doc); err != nil {
			log.Fatalf("Failed to unmarshal document: %v", err)
		}

		// Display full document details
		fmt.Printf("=== Document Details ===\n\n")
		fmt.Printf("Object ID: %s\n", doc.ObjectID)
		fmt.Printf("URL: %s\n", doc.URL)
		if doc.Anchor != "" {
			fmt.Printf("Anchor: %s\n", doc.Anchor)
		}
		fmt.Println()

		fmt.Printf("Hierarchy:\n")
		if doc.HierarchyLvl0 != nil {
			fmt.Printf("  Level 0: %s\n", *doc.HierarchyLvl0)
		}
		if doc.HierarchyLvl1 != nil {
			fmt.Printf("  Level 1: %s\n", *doc.HierarchyLvl1)
		}
		if doc.HierarchyLvl2 != nil {
			fmt.Printf("  Level 2: %s\n", *doc.HierarchyLvl2)
		}
		if doc.HierarchyLvl3 != nil {
			fmt.Printf("  Level 3: %s\n", *doc.HierarchyLvl3)
		}
		if doc.HierarchyLvl4 != nil {
			fmt.Printf("  Level 4: %s\n", *doc.HierarchyLvl4)
		}
		if doc.HierarchyLvl5 != nil {
			fmt.Printf("  Level 5: %s\n", *doc.HierarchyLvl5)
		}
		if doc.HierarchyLvl6 != nil {
			fmt.Printf("  Level 6: %s\n", *doc.HierarchyLvl6)
		}
		fmt.Println()

		fmt.Printf("Radio Hierarchy:\n")
		if doc.HierarchyRadioLvl0 != nil {
			fmt.Printf("  Radio Level 0: %s\n", *doc.HierarchyRadioLvl0)
		}
		if doc.HierarchyRadioLvl1 != nil {
			fmt.Printf("  Radio Level 1: %s\n", *doc.HierarchyRadioLvl1)
		}
		if doc.HierarchyRadioLvl2 != nil {
			fmt.Printf("  Radio Level 2: %s\n", *doc.HierarchyRadioLvl2)
		}
		if doc.HierarchyRadioLvl3 != nil {
			fmt.Printf("  Radio Level 3: %s\n", *doc.HierarchyRadioLvl3)
		}
		if doc.HierarchyRadioLvl4 != nil {
			fmt.Printf("  Radio Level 4: %s\n", *doc.HierarchyRadioLvl4)
		}
		if doc.HierarchyRadioLvl5 != nil {
			fmt.Printf("  Radio Level 5: %s\n", *doc.HierarchyRadioLvl5)
		}
		fmt.Println()

		fmt.Printf("Content:\n")
		if doc.Content != nil {
			fmt.Printf("%s\n", *doc.Content)
		} else {
			fmt.Printf("(no content)\n")
		}
	},
}
