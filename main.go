package main

import (
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

type Link struct {
	url   string
	title string
}

// TODO: Write your function here
func CrawlLink(url string) []Link {
	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error")
		panic(err.Error)
	}
	defer response.Body.Close()

	documment, err := html.Parse(response.Body)
	if err != nil {
		fmt.Println("Error")
		panic(err.Error)
	}

	nodes := getNode(documment)

	var links []Link
	for _, node := range nodes {
		var link Link
		for _, attribute := range node.Attr {
			if attribute.Key == "href" {
				link.url = attribute.Val
			}
		}
		link.title = getText(node)

		links = append(links, link)
	}

	return links
}

func getNode(node *html.Node) []*html.Node {
	if node.Type == html.ElementNode && node.Data == "a" {
		return []*html.Node{
			node,
		}
	}

	var nodes []*html.Node
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		nodes = append(nodes, getNode(child)...)
	}
	return nodes
}

func getText(node *html.Node) string {
	if node.Type == html.TextNode {
		return node.Data
	}

	var text string
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		text += getText(child)
	}
	return strings.Join(strings.Fields(text), " ")
}
