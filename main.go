package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"

	"golang.org/x/net/html"
)

var url string

func init() {
	flag.StringVar(&url, `url`, ``, `URL to save locally`)
	flag.StringVar(&url, `u`, ``, `URL to save locally`)
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
	raw, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("Error: unable to read response body: %+v", err)
		return
	}

	doc, err := html.Parse(res.Body)
	if err != nil {
		fmt.Printf("Error: unable to parse response body: %+v, %v", err, doc)
		return
	}

	fmt.Printf(`res: %+v`, string(raw))

}
