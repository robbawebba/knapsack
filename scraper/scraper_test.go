package scraper

import (
	"os"
	"testing"
)

func TestScrapeArticle(t *testing.T) {
	testMaterialURI := []string{"../content/monolith-vs-microservice-vs-serverless-the-real-winner-the-developer-8aae6042fb48.html"}
	scraper := new(Scraper)
	for _, uri := range testMaterialURI {
		testFile, err := os.Open(uri)
		if err != nil {
			t.Errorf("Failed to open test file at %s", uri)
			continue
		}
		err = scraper.ScrapeArticle(testFile)
		if err != nil {
			t.Error(err)
		}

	}
}
