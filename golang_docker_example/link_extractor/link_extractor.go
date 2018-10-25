package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"golang.org/x/net/html"
)

const (
	httpPrefix = "http://"
)

type ServerResponse struct {
	Page  string   `json:"page"`
	Links []string `json:"links"`
	// msg   string   `json:"msg"`
}
type ServerErrorResponse struct {
	Msg string `json:"msg"`
}
type serverCoords struct {
	address string
	port    string
}

func main() {
	serverCoords, err := getServerCoords()
	must(err)

	r := mux.NewRouter()
	r.HandleFunc("/extract", extractor)
	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(serverCoords.address+":"+serverCoords.port, nil))
}

func getServerCoords() (*serverCoords, error) {
	if len(os.Args[1]) == 0 {
		return nil, errors.New("Http address cannot be empty")
	}
	_, err := strconv.Atoi(os.Args[2])
	if err != nil {
		return nil, errors.New("Port should be a numeric value")
	}
	return &serverCoords{address: os.Args[1], port: os.Args[2]}, nil
}

func must(err error) {
	if err != nil {
		log.Fatal("Incorrect server settings", err)
	}
}

func extractor(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("from")

	page, err := parse(httpPrefix + url)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(err)
		return
	}

	links := extractLinks(nil, page)
	resp := ServerResponse{Page: url, Links: links}

	if len(links) == 0 {
		fmt.Println("2")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(&resp)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&resp)
}

func parse(url string) (*html.Node, error) {
	r, err := http.Get(url)
	if err != nil {
		return nil, errors.New("Page cannot be reached")
	}

	b, err := html.Parse(r.Body)
	if err != nil {
		return nil, errors.New("Page cannot be parsed")
	}

	return b, nil

}

func extractLinks(links []string, n *html.Node) []string {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key == "href" {

				links = append(links, a.Val)
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		links = extractLinks(links, c)
	}

	return links

}
func extractLinksConcurrently(links []string, n *html.Node, ch chan<- string) {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key == "href" {
				ch <- a.Val
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		extractLinksConcurrently(links, c, ch)
	}

}
