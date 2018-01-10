package main

import (
	"flag"
	"fmt"
	"net/http"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

var url string

func init() {
	flag.StringVar(&url, `url`, ``, `URL to save locally`)
	flag.StringVar(&url, `u`, ``, `URL to save locally`)
}

func findAnchors(n *html.Node) int {
	anc := 0
	if n.Type == html.ElementNode && n.DataAtom == atom.A {
		anc++
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		anc += findAnchors(c)
	}
	return anc
}

func findAnchorsWithTokenizer(t *html.Tokenizer) {
	count := 0
	for {
		tt := t.Next()
		switch tt {
		case html.ErrorToken:
			return
		case html.StartTagToken, html.EndTagToken:
			tn := t.Token()
			if tn.DataAtom == atom.A {
				if tt == html.StartTagToken {
					count++
				}
			}
		}
	}
}

func main() {
	flag.Parse()
	if url == `` {
		return
	}
	res, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error requesting URL %s: %+v", url, err)
		return
	}
	defer res.Body.Close()
	if err != nil {
		fmt.Printf("Error: unable to read response body: %+v", err)
		return
	}

	doc, err := html.Parse(res.Body)
	if err != nil {
		fmt.Printf("Error: unable to parse response body: %+v, %v", err, doc)
		return
	}
	count := findAnchors(doc)

	fmt.Printf(`count: %d`, count)

}
