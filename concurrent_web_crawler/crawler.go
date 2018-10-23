package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"golang.org/x/net/html"
)

type link struct {
	depth int
	value string
}

func crawl(url link) []link {
	fmt.Println(url)
	list, err := extract(url.value)
	if err != nil {
		log.Print(err)
	}
	var result []link
	for _, e := range list {
		result = append(result, link{depth: url.depth + 1, value: e})
	}
	return result
}

func mapToLink(args []string, d int) []link {
	var result []link
	for _, a := range args {
		result = append(result, link{depth: d, value: a})
	}
	return result
}

func filterUpToDepth(args []link, d int) []link {
	var result []link
	for _, a := range args {
		if a.depth < d {
			result = append(result, a)
		}
	}
	return result
}

func main() {
	worklist := make(chan []link)
	unseenLinks := make(chan link)
	seen := make(map[string]bool)

	go func() { worklist <- mapToLink(os.Args[2:], 0) }() //1 - depth

	for i := 0; i < 20; i++ {
		go func() {
			for link := range unseenLinks {
				foundLinks := crawl(link)
				go func() { worklist <- foundLinks }()
			}
		}()
	}

	for list := range worklist {
		for _, link := range list {
			d, err := strconv.Atoi(os.Args[1])
			if err == nil {
				if !seen[link.value] && link.depth < d {
					seen[link.value] = true
					unseenLinks <- link
				}
			}
		}
	}

}

func extract(url string) ([]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("Download %s: %s", url, resp.Status)
	}

	doc, err := html.Parse(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("Parse %s as HTML: %v", url, err)
	}

	var links []string
	visitNode := func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key != "href" {
					continue
				}
				link, err := resp.Request.URL.Parse(a.Val)
				if err != nil {
					continue
				}
				links = append(links, link.String())
			}
		}
	}
	forEachNode(doc, visitNode, nil)
	return links, nil
}

func forEachNode(n *html.Node, pre, post func(n *html.Node)) {
	if pre != nil {
		pre(n)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		forEachNode(c, pre, post)
	}
	if post != nil {
		post(n)
	}
}
