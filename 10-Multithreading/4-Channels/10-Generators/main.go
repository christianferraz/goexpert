package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"sync"
	"time"
)

var titleRegex = regexp.MustCompile(`<title>(.*?)</title>`)

func titulo(url string, c chan string) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	html, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	titulo := titleRegex.FindStringSubmatch(string(html))

	if len(titulo) > 1 {
		c <- titulo[1]
	} else {
		c <- fmt.Sprintf("nao há titulo em %s", url)
	}
}

func main() {
	var wg sync.WaitGroup
	c := make(chan string)

	urls := []string{
		"https://www.cod3r.com.br",
		"https://www.google.com",
		"https://www.amazon.com",
		"https://www.youtube.com",
		"https://troia2025.sbop.org.br",
	}

	for _, url := range urls {
		wg.Go(func() {
			titulo(url, c)
		})
	}

	go func() {
		wg.Wait()
		close(c)
	}()

	for titulo := range c {
		fmt.Println(titulo)
	}
}
