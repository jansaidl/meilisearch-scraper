package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var RootCmd = &cobra.Command{
	Use:   "meilisearch-scraper",
	Short: "Meilisearch documentation scraper",
	Long: `A powerful CLI tool for scraping documentation websites and indexing them in Meilisearch.
It parses sitemaps, extracts content using CSS selectors, and creates searchable documents.`,
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	RootCmd.PersistentFlags().String("config", "", "config file path (default: config.json)")
	RootCmd.PersistentFlags().String("meilisearch-url", "", "Meilisearch server URL (env: MEILISEARCH_HOST_URL)")
	RootCmd.PersistentFlags().String("meilisearch-key", "", "Meilisearch API key (env: MEILISEARCH_API_KEY)")
	RootCmd.PersistentFlags().String("index", "docs", "Meilisearch index name (env: MEILISEARCH_INDEX)")

	viper.BindPFlag("config", RootCmd.PersistentFlags().Lookup("config"))
	viper.BindPFlag("meilisearch.url", RootCmd.PersistentFlags().Lookup("meilisearch-url"))
	viper.BindPFlag("meilisearch.key", RootCmd.PersistentFlags().Lookup("meilisearch-key"))
	viper.BindPFlag("meilisearch.index", RootCmd.PersistentFlags().Lookup("index"))

	// Bind environment variables
	viper.BindEnv("meilisearch.url", "MEILISEARCH_HOST_URL")
	viper.BindEnv("meilisearch.key", "MEILISEARCH_API_KEY")
	viper.BindEnv("meilisearch.index", "MEILISEARCH_INDEX")
	viper.BindEnv("sitemap.url", "SITEMAP_URL")
	viper.BindEnv("config", "CONFIG_PATH")

	// Add all commands
	RootCmd.AddCommand(runCmd)
	RootCmd.AddCommand(dryRunCmd)
	RootCmd.AddCommand(deleteCmd)
	RootCmd.AddCommand(testCmd)
	RootCmd.AddCommand(inspectCmd)
	RootCmd.AddCommand(statsCmd)
	RootCmd.AddCommand(listCmd)
	RootCmd.AddCommand(detailCmd)
	RootCmd.AddCommand(searchCmd)
}

func initConfig() {
	configFile := viper.GetString("config")
	if configFile == "" {
		configFile = "config.json"
	}
	viper.SetConfigFile(configFile)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		log.Printf("Using config file: %s", viper.ConfigFileUsed())
	}
}
