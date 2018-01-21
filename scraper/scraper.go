package scraper

import (
	"errors"
	"io"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type Scraper struct {
}

// ScrapeArticle searches the provided reader for information about an article: author, title, created, body
func (s *Scraper) ScrapeArticle(r io.Reader) error {
	doc, err := html.Parse(r)
	if err != nil {
		return err
	}

	//first find body and head
	head := findNode(doc, atom.Head)
	if head == nil {
		return errors.New("Could not find a head tag, perhaps this is invalid html")
	}
	body := findNode(doc, atom.Body)
	if body == nil {
		return errors.New("Could not find a body tag, perhaps this is invalid html")
	}
	return nil
}

func findNode(root *html.Node, seekTag atom.Atom) *html.Node {
	var traverse func(*html.Node)
	var foundNode *html.Node

	traverse = func(node *html.Node) {
		if node.Type == html.ElementNode && node.DataAtom == seekTag {
			foundNode = node
			return
		}
		if node.Type == html.DoctypeNode && root.NextSibling != nil {
			findNode(root.NextSibling, seekTag)
		}
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}
	traverse(root)
	return foundNode
}
