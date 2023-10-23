package main

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"sync"
)

func titulo(url string, wg *sync.WaitGroup, c chan string) {
	defer wg.Done()
	resp, err := http.Get(url)
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

	r, err := regexp.Compile(`<title>(.*?)</title>`)
	if err != nil {
		fmt.Println(err)
		return
	}

	titulo := r.FindStringSubmatch(string(html))
	if len(titulo) > 1 {
		c <- titulo[1]
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
	}

	for _, url := range urls {
		wg.Add(1)
		go titulo(url, &wg, c)
	}

	go func() {
		wg.Wait()
		close(c)
	}()

	for titulo := range c {
		fmt.Println(titulo)
	}
}
