package main

import (
	"flag"
	"fmt"
	"strings"
	"time"
	//"html"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

//var a = foo()

//Simple slice based queue
type Queue []string

func (q *Queue) Add(str string) {
	*q = append(*q, str)
}

func (q *Queue) Remove() (string, bool) {
	if len(*q) == 0 {
		return "", false
	} else {
		str := (*q)[0]
		*q = (*q)[1:]
		return str, true
	}
}

type StatDiv struct {
	Value      string
	TextLength int
	DotCount   int
	Depth      int
}

var (
	WORKERS  int    = 5
	TIME_OUT int    = 100
	HTTPS    bool   = true
	HOST     string = "meduza.io"
)

func grabber() <-chan string {
	stringChan := make(chan string)
	for i := 0; i < WORKERS; i++ {
		go func() {
			for {
				//do smth
				time.Sleep(time.Duration(TIME_OUT) * time.Millisecond)
			}
		}()
	}
	fmt.Println(WORKERS, " have been launched")
	return stringChan
}

func depth(node *html.Node) int {
	if node.Parent != nil {
		return 1 + depth(node.Parent)
	}

	return 0
}

func getDoc(link string, m map[string]int, links *Queue) {
	doc, err := goquery.NewDocument(link)
	if err != nil {
		fmt.Println("Error: ")
		fmt.Println(err)
	}

	var max, maxValue int = 0, 0

	//Searching main article
	mapDivs := make(map[int]StatDiv)
	mapClass := make(map[string]int)
	divs := doc.Find("div")
	divs.Each(func(i int, s *goquery.Selection) {

		length := len(s.Text())
		count := strings.Count(s.Text(), ".")

		stat := mapDivs[i]
		stat.DotCount = count
		stat.TextLength = length
		mapDivs[i] = stat

		node := s.Nodes[0]

		fmt.Printf("ID: %d; Nodes: %d; depth - %d; count - %d; length - %d; attributes: count - %d; ", i, len(s.Nodes), depth(node), count, length, len(node.Attr))

		for _, item := range node.Attr {
			if item.Key == "class" {
				//			if countClass, exists := mapClass[item.Val]; exists {
				//				mapClass[item.Val]
				//			}
				mapClass[item.Val]++
				fmt.Printf("value -%s, ", item.Val)
			}
		}
		fmt.Printf("\n")

		if count > 2 && length/count < 300 {
			fmt.Println(s.Text())
			attr := node.Attr[0]
			fmt.Printf("Namespace: %s; key:%s; value:%s \n", attr.Namespace, attr.Key, attr.Val)
		}

	})
	fmt.Println(max, maxValue)
	/*
		for key, val := range mapDivs {
			fmt.Printf("ID: %d; TextCount: %d; DotCount: %d \n", key, val.TextLength, val.DotCount)
		}

	*/

	//Searching links
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		//		doc.Find(".sidebar-reviews article .content-block").Each(func(i int, s *goquery.Selection) {

		if href, ok := s.Attr("href"); ok {
			count := m[href]
			if href[0] == 47 {
				if count == 0 {
					links.Add(constructURL(href))
				}
				m[href] = count + 1
			}
			//fmt.Printf("%d - %d - %s - %s - %d\n", i, href[0], s.Text(), href, count)
		}
	})
}

func constructURL(link string) string {
	var url string

	if HTTPS {
		url = "https://"
	} else {
		url = "http://"
	}

	url = url + HOST + link
	return url
}

func main() {
	flag.IntVar(&WORKERS, "w", WORKERS, "workers")
	flag.IntVar(&TIME_OUT, "t", TIME_OUT, "time out")
	flag.BoolVar(&HTTPS, "s", HTTPS, "use https")
	flag.StringVar(&HOST, "h", HOST, "host")
	flag.Parse()

	m := make(map[string]int)

	links := make(Queue, 0)

	//getDoc(constructURL(""), m, &links)
	fmt.Println(links)
	fmt.Println("Links for downloading: ", len(links))
	fmt.Println("Total links: ", len(m))
	/*
		var i int = 0
			for next, exists := links.Remove(); exists; {
				getDoc(next, m, &links)
				fmt.Println("Next: ", next)
				fmt.Println(links)
				fmt.Println("Links for downloading: ", len(links))
				fmt.Println("Total links: ", len(m))
				i++
				if i > 5 {
					break
				}
				next, exists = links.Remove()
			}
	*/

	getDoc("https://meduza.io/news/2018/03/11/v-kitae-otmenili-ogranichenie-sroka-prebyvaniya-u-vlasti-dlya-lidera-strany", m, &links)

}
