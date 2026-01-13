package src

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func FetchSitemap(url string) (*Sitemap, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch sitemap: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var sitemap Sitemap
	if err := xml.Unmarshal(body, &sitemap); err != nil {
		return nil, fmt.Errorf("failed to parse sitemap XML: %w", err)
	}

	return &sitemap, nil
}

func ScrapePage(pageURL string, config *Config) ([]Document, error) {
	// Always use .html extension to get static content instead of JS-rendered version
	fetchURL := pageURL
	if !strings.HasSuffix(pageURL, ".html") {
		fetchURL = pageURL + ".html"
	}

	resp, err := http.Get(fetchURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch page: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	goDoc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	var documents []Document

	// Extract global lvl0 if configured
	var lvl0Value *string
	if config.Selectors.Lvl0.Global {
		text := goDoc.Find(config.Selectors.Lvl0.Selector).First().Text()
		text = strings.TrimSpace(text)
		if text == "" && config.Selectors.Lvl0.DefaultValue != "" {
			text = config.Selectors.Lvl0.DefaultValue
		}
		if text != "" {
			lvl0Value = &text
		}
	}

	// Extract hierarchy levels
	var lvl1Value, lvl2Value, lvl3Value, lvl4Value, lvl5Value, lvl6Value *string

	if config.Selectors.Lvl1 != "" {
		text := goDoc.Find(config.Selectors.Lvl1).First().Text()
		text = strings.TrimSpace(text)
		if text != "" {
			lvl1Value = &text
		}

		// Process each heading section
		goDoc.Find(config.Selectors.Lvl2).Each(func(i int, s *goquery.Selection) {
			anchor, _ := s.Attr("id")
			if anchor == "" {
				anchor = fmt.Sprintf("section_%d", i)
			}

			lvl2Text := strings.TrimSpace(s.Text())
			if lvl2Text != "" {
				lvl2Value = &lvl2Text
			}

			// Extract content for this section
			var contentParts []string

			// Get text content following this heading
			s.NextUntil(config.Selectors.Lvl2).Each(func(j int, content *goquery.Selection) {
				if config.Selectors.Text != "" {
					content.Find(config.Selectors.Text).Each(func(k int, textNode *goquery.Selection) {
						text := strings.TrimSpace(textNode.Text())
						if text != "" {
							contentParts = append(contentParts, text)
						}
					})
				}
			})

			contentText := strings.Join(contentParts, " ")
			var contentPtr *string
			if contentText != "" {
				contentPtr = &contentText
			}

			fullURL := pageURL
			if anchor != "" {
				fullURL = fmt.Sprintf("%s#%s", pageURL, anchor)
			}

			// Generate objectID
			hash := sha256.Sum256([]byte(fullURL))
			objectID := hex.EncodeToString(hash[:])

			doc := Document{
				Anchor:             anchor,
				Content:            contentPtr,
				URL:                fullURL,
				ObjectID:           objectID,
				HierarchyLvl0:      lvl0Value,
				HierarchyLvl1:      lvl1Value,
				HierarchyLvl2:      lvl2Value,
				HierarchyLvl3:      lvl3Value,
				HierarchyLvl4:      lvl4Value,
				HierarchyLvl5:      lvl5Value,
				HierarchyLvl6:      lvl6Value,
				HierarchyRadioLvl0: nil,
				HierarchyRadioLvl1: lvl1Value,
				HierarchyRadioLvl2: lvl2Value,
				HierarchyRadioLvl3: lvl3Value,
				HierarchyRadioLvl4: lvl4Value,
				HierarchyRadioLvl5: lvl5Value,
			}

			documents = append(documents, doc)
		})

		// If no lvl2 headings found, create at least one document for the page
		if len(documents) == 0 {
			var contentParts []string
			if config.Selectors.Text != "" {
				goDoc.Find(config.Selectors.Text).Each(func(i int, s *goquery.Selection) {
					text := strings.TrimSpace(s.Text())
					if text != "" {
						contentParts = append(contentParts, text)
					}
				})
			}

			contentText := strings.Join(contentParts, " ")
			var contentPtr *string
			if contentText != "" {
				contentPtr = &contentText
			}

			hash := sha256.Sum256([]byte(pageURL))
			objectID := hex.EncodeToString(hash[:])

			doc := Document{
				Anchor:             "",
				Content:            contentPtr,
				URL:                pageURL,
				ObjectID:           objectID,
				HierarchyLvl0:      lvl0Value,
				HierarchyLvl1:      lvl1Value,
				HierarchyLvl2:      nil,
				HierarchyLvl3:      nil,
				HierarchyLvl4:      nil,
				HierarchyLvl5:      nil,
				HierarchyLvl6:      nil,
				HierarchyRadioLvl0: nil,
				HierarchyRadioLvl1: lvl1Value,
				HierarchyRadioLvl2: nil,
				HierarchyRadioLvl3: nil,
				HierarchyRadioLvl4: nil,
				HierarchyRadioLvl5: nil,
			}

			documents = append(documents, doc)
		}

	}
	return documents, nil
}
