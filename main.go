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

func findHeadersWithTokenizer(t *html.Tokenizer, foundHeaders chan []html.Token) {
	for {
		tt := t.Next()
		switch tt {
		case html.ErrorToken:
			foundHeaders <- []html.Token{t.Token()}
			return
		case html.StartTagToken:
			tn := t.Token()
			switch tn.DataAtom {
			case atom.H1, atom.H2, atom.H3, atom.H4, atom.H5, atom.H6:
				if tt == html.StartTagToken {
					headerTree := []html.Token{tn}
					// search for the rest of the tokens of this header node
					for {
						t.Next()
						n := t.Token()
						headerTree = append(headerTree, n)
						if n.DataAtom == tn.DataAtom && n.Type == html.EndTagToken {
							break
						}
					}
					foundHeaders <- headerTree

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

	t := html.NewTokenizer(res.Body)
	if err != nil {
		fmt.Printf("Error: unable to parse response body: %+v, %v", err, t)
		return
	}

	headers := make(chan []html.Token)
	go findHeadersWithTokenizer(t, headers)

loop:
	for {
		select {
		case h := <-headers:
			if h[0].Type == html.ErrorToken {
				break loop
			}
			fmt.Printf("Found header %+v\n", h)
		}
	}

}
