# Meilisearch Documentation Scraper

Powerful CLI tool for scraping documentation websites and indexing them in Meilisearch. It parses sitemaps, extracts content using CSS selectors, and creates searchable documents with hierarchical structure.

## Features

- **Sitemap-based scraping** - Automatically discovers and scrapes all URLs from XML sitemaps
- **Configurable CSS selectors** - Extract content using customizable CSS selectors for different hierarchy levels
- **Meilisearch integration** - Direct upload to Meilisearch for instant full-text search
- **Document management** - List, search, and view detailed information about indexed documents
- **Testing tools** - Dry-run mode and single-page testing for configuration validation
- **HTML inspection** - Analyze page structure to find the right CSS selectors

## Installation

```bash
go install github.com/jansaidl/meilisearch-scraper/cmd@latest
```

Or build from source:

```bash
git clone https://github.com/jansaidl/meilisearch-scraper
cd meilisearch-scraper
go build -o meilisearch-scraper ./cmd
```

## Configuration

### Environment Variables

```bash
export MEILISEARCH_HOST_URL="http://localhost:7700"
export MEILISEARCH_API_KEY="your-api-key"
export MEILISEARCH_INDEX="docs"
export SITEMAP_URL="https://docs.example.com/sitemap.xml"
```

### Config File (config.json)

The tool uses a JSON configuration file to define CSS selectors for content extraction:

```json
{
  "selectors": {
    "lvl0": {
      "selector": "nav.sidebar .active",
      "global": true,
      "default_value": "Documentation"
    },
    "lvl1": "article h1",
    "lvl2": "article h2",
    "lvl3": "article h3",
    "lvl4": "article h4",
    "lvl5": "article h5",
    "lvl6": "article h6",
    "text": "article p, article li"
  }
}
```

## Commands

### `run` - Scrape and Upload

Scrapes all URLs from a sitemap and uploads documents to Meilisearch.

```bash
# Basic usage
meilisearch-scraper run https://docs.example.com/sitemap.xml

# With environment variable
export SITEMAP_URL=https://docs.example.com/sitemap.xml
meilisearch-scraper run

# Limit number of URLs to process
meilisearch-scraper run --limit 10

# Custom config file
meilisearch-scraper run --config my-config.json
```

**Flags:**
- `--limit` - Limit number of URLs to process (0 = no limit)
- `--config` - Config file path (default: config.json)
- `--index` - Meilisearch index name (default: docs)

---

### `dry-run` - Test Without Upload

Scrapes documentation and saves results to `data.json` without uploading to Meilisearch. Useful for testing and debugging configurations.

```bash
# Test scraping configuration
meilisearch-scraper dry-run https://docs.example.com/sitemap.xml

# Test with limited URLs
meilisearch-scraper dry-run https://docs.example.com/sitemap.xml --limit 5
```

**Flags:**
- `--limit` - Limit number of URLs to process

---

### `test` - Test Single URL

Test the scraping configuration on a single URL and view extracted documents in JSON format.

```bash
# Test single page
meilisearch-scraper test https://docs.example.com/getting-started

# Test with custom config
meilisearch-scraper test https://docs.example.com/api --config my-config.json
```

---

### `search` - Search Documents

Search for documents in the Meilisearch index using full-text search.

```bash
# Basic search
meilisearch-scraper search "installation"

# Search in specific index
meilisearch-scraper search "api" --index my-docs

# Limit number of results
meilisearch-scraper search "configuration" --limit 5
```

**Flags:**
- `--limit` - Maximum number of results to return (default: 10)
- `--index` - Meilisearch index name

---

### `list` - List Documents

List all documents in the Meilisearch index with their titles and URLs.

```bash
# List all documents
meilisearch-scraper list

# List from specific index
meilisearch-scraper list --index my-docs

# Limit results
meilisearch-scraper list --limit 20
```

**Flags:**
- `--limit` - Maximum number of documents to list (default: 100)

---

### `detail` - Show Document Details

Display complete information about a specific document using its objectID.

```bash
# Show document by ID
meilisearch-scraper detail abc123def456...

# From specific index
meilisearch-scraper detail abc123def456... --index my-docs
```

---

### `inspect` - Inspect HTML Elements

Inspect HTML elements at a URL using CSS selectors. Helps you understand page structure and test selectors before scraping.

```bash
# Inspect h1 elements
meilisearch-scraper inspect https://docs.example.com/page "h1"

# Inspect navigation
meilisearch-scraper inspect https://docs.example.com/page "nav.sidebar"

# Inspect with class selector
meilisearch-scraper inspect https://docs.example.com/page ".content h2"
```

---

### `stats` - Index Statistics

Display statistics about the Meilisearch index including document count and field distribution.

```bash
# Show stats for default index
meilisearch-scraper stats

# Show stats for specific index
meilisearch-scraper stats --index my-docs
```

---

### `delete` - Delete All Documents

Delete all documents from the specified Meilisearch index.

```bash
# Delete from default index
meilisearch-scraper delete

# Delete from specific index
meilisearch-scraper delete --index my-docs
```

## Document Structure

Each scraped document contains:

```json
{
  "objectID": "unique-document-id",
  "url": "https://docs.example.com/page",
  "anchor": "section-name",
  "hierarchy_lvl0": "Main Section",
  "hierarchy_lvl1": "Subsection",
  "hierarchy_lvl2": "Heading",
  "hierarchy_lvl3": "Subheading",
  "hierarchy_lvl4": "...",
  "hierarchy_lvl5": "...",
  "hierarchy_lvl6": "...",
  "hierarchy_radio_lvl0": "...",
  "hierarchy_radio_lvl1": "...",
  "content": "Extracted text content"
}
```

## Global Flags

These flags can be used with any command:

- `--config` - Config file path (default: config.json)
- `--meilisearch-url` - Meilisearch server URL
- `--meilisearch-key` - Meilisearch API key
- `--index` - Meilisearch index name (default: docs)

## Workflow Example

1. **Inspect a page** to find the right CSS selectors:
   ```bash
   meilisearch-scraper inspect https://docs.example.com/page "article h1"
   ```

2. **Create configuration file** with appropriate selectors (`config.json`)

3. **Test single page** to verify extraction:
   ```bash
   meilisearch-scraper test https://docs.example.com/page
   ```

4. **Dry run** to test full sitemap without uploading:
   ```bash
   meilisearch-scraper dry-run https://docs.example.com/sitemap.xml --limit 5
   ```

5. **Run full scraping** and upload to Meilisearch:
   ```bash
   meilisearch-scraper run https://docs.example.com/sitemap.xml
   ```

6. **Search documents** to verify indexing:
   ```bash
   meilisearch-scraper search "getting started"
   ```

7. **View statistics** to check index health:
   ```bash
   meilisearch-scraper stats
   ```

## Requirements

- Go 1.21 or higher
- Meilisearch server running and accessible
- Valid Meilisearch API key

## Dependencies

- [meilisearch-go](https://github.com/meilisearch/meilisearch-go) - Meilisearch Go client
- [cobra](https://github.com/spf13/cobra) - CLI framework
- [viper](https://github.com/spf13/viper) - Configuration management
- [goquery](https://github.com/PuerkitoBio/goquery) - HTML parsing

## License

MIT License

