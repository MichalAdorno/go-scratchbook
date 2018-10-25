# Concurrent web crawler
## (based on the Exercise 8.6 from Donovan, Kernighan, The Go Programming Language)

Goroutines are used to traverse links concurrently until they find out that all links have reached given depth counting from the initial links).

The program crawls the web given:
* a depth param (any int > 0)
* initial http links (which are the starting points for the crawler).

Usage:
```
    go clean
    go build
    ./crawler [depth] [list of space-separated http addresses]
```
For example:
```
    ./crawler 3 http://wikipedia.com 1> links.txt 2> errors.txt
```




