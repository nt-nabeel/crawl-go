package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

func main() {
	ssr := CrawlLink("https://ssr-assessment-miqdadyyy.vercel.app/")
	csr := CrawlLink("https://csr-assessment-miqdadyyy.vercel.app/")

	fmt.Println("")
	fmt.Println(ssr)
	fmt.Println("")
	fmt.Println("")
	fmt.Println(csr)
	fmt.Println("")
}

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
		for _, attribute := range node.Attr {
			if attribute.Key == "href" {
				var link Link
				link.url = attribute.Val
				link.title = getText(node)
				links = append(links, link)
			}
			if attribute.Key == "src" {
				nameSlice := strings.Split(attribute.Val, ".")
				isApp := strings.HasSuffix(nameSlice[0], "app")
				if isApp {
					links = getLinkFromJS(url+attribute.Val, links)
				}
			}
		}
	}

	return links
}

func getLinkFromJS(url string, links []Link) []Link {
	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error")
		panic(err.Error)
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error")
		panic(err.Error)
	}
	data := strings.Split(string(body), ";")
	for _, c := range data {
		if len(strings.Split(c, "links=")) == 2 {
			_links := strings.Split(c, "links=")[1]
			_links = _links[2:(len(_links) - 2)]
			linkObj := strings.Split(_links, "},{")
			for _, _link := range linkObj {
				_link = _link[5:]
				_data := strings.Split(_link, `",title:"`)

				var link Link
				link.url = _data[0]
				link.title = _data[1][:(len(_data[1]) - 1)]
				links = append(links, link)
			}
		}
	}
	return links
}

func getNode(node *html.Node) []*html.Node {
	if node.Type == html.ElementNode && (node.Data == "a" || node.Data == "script") {
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
