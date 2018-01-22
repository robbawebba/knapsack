package scraper

import (
	"errors"
	"fmt"
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
	head := findFirst(doc, atom.Head)
	if head == nil {
		return errors.New("Could not find a head tag, perhaps this is invalid html")
	}
	body := findFirst(doc, atom.Body)
	if body == nil {
		return errors.New("Could not find a body tag, perhaps this is invalid html")
	}

	headers := findAll(doc, atom.H1, atom.H2, atom.H3, atom.H4, atom.H5)
	if headers == nil {
		return errors.New("Could not find any header tags")
	}
	fmt.Println(" FOUND HEADERS")
	for _, h := range headers {
		fmt.Printf("%+v\n", h)
	}
	return nil
}

func findFirst(root *html.Node, seekTag atom.Atom) *html.Node {
	var traverse func(*html.Node)
	var foundNode *html.Node

	traverse = func(node *html.Node) {
		if node.Type == html.ElementNode && node.DataAtom == seekTag {
			foundNode = node
			return
		}
		if node.Type == html.DoctypeNode && root.NextSibling != nil {
			// findNode(root.NextSibling, seekTag)
		}
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}
	traverse(root)
	return foundNode
}

func findAll(root *html.Node, seekTags ...atom.Atom) []*html.Node {
	var traverse func(*html.Node)
	var foundNodes []*html.Node

	traverse = func(node *html.Node) {
		for _, tag := range seekTags {
			if node.Type == html.ElementNode && node.DataAtom == tag {
				foundNodes = append(foundNodes, node)
			}
		}
		if node.Type == html.DoctypeNode && root.NextSibling != nil {
			// findNod(root.NextSibling, seekTag)
		}
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}
	traverse(root)
	return foundNodes
}
