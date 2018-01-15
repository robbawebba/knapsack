package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

var outPath string

func init() {
	flag.StringVar(&outPath, `out`, `nooooo`, `Save the HTML body of the response to a file`)
	flag.StringVar(&outPath, `o`, `nooooo`, `Save the HTML body of the response to a file`)
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
	url := flag.Arg(0)

	if url == `` {
		fmt.Println(`NO URL`)
		return
	}
	res, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error requesting URL %s: %+v", url, err)
		return
	}

	defer res.Body.Close()

	if outPath != `` {
		fmt.Println(outPath)
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Printf("Error: unable to read response body: %+v", err)
			return
		}

		file, err := os.Create(outPath)
		if err != nil {
			fmt.Printf("Error creating new file: %+v", err)
			return
		}
		defer file.Close()
		_, err = file.Write(body)
		if err != nil {
			fmt.Printf("Error while writing body to file: %+v", err)
			return
		}
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
