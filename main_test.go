package main

import (
	"fmt"
	"net/http"
	"testing"

	"golang.org/x/net/html"
)

func BenchmarkFindAnchors(b *testing.B) {
	res, err := http.Get("https://medium.freecodecamp.org/the-crazy-history-of-the-100daysofcode-challenge-and-why-you-should-try-it-for-2018-6c89a76e298d")
	if err != nil {
		fmt.Printf("Error requesting URL: %+v", err)
		return
	}
	defer res.Body.Close()
	if err != nil {
		fmt.Printf("Error: unable to read response body: %+v", err)
		return
	}

	dock, err := html.Parse(res.Body)
	if err != nil {
		fmt.Printf("Error: unable to parse response body: %+v, %v", err, dock)
		return
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = findAnchors(dock)
	}
}

func BenchmarkFindHeadersWithTokenizer(b *testing.B) {
	res, err := http.Get("https://medium.freecodecamp.org/the-crazy-history-of-the-100daysofcode-challenge-and-why-you-should-try-it-for-2018-6c89a76e298d")
	if err != nil {
		fmt.Printf("Error requesting URL: %+v", err)
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
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
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
}
