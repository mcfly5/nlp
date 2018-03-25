package main

import (
	"flag"
	"fmt"
	"hash/crc32"
	"io/ioutil"
	"os"
	"strconv"
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
	WORKERS   int    = 5
	TIME_OUT  int    = 100
	HTTPS     bool   = true
	HOST      string = "meduza.io"
	DATA_DIR  string = "data/"
	MAX_LINKS int    = 50
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

func getLinks(doc *goquery.Document, m map[string]int, links *Queue) {

	//Searching links
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		if href, ok := s.Attr("href"); ok {
			count := m[href]
			if href[0] == 47 {
				if count == 0 {
					links.Add(constructURL(href))
					fmt.Printf("New link found : %d - %d - %s - %s \n", i, href[0], s.Text(), href)
				}
				m[href] = count + 1
			}
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

func getArticle(doc *goquery.Document) {

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

}

//Read goquery doc from file
func readDocumentFromFile(fileName string) (doc *goquery.Document, err error) {

	bs, err := ioutil.ReadFile(fileName)
	if err != nil {
		return doc, err
	}

	strFromFile := string(bs)
	fmt.Println(strFromFile)

	ior := strings.NewReader(strFromFile)

	doc, err = goquery.NewDocumentFromReader(ior)
	if err != nil {
		return doc, err
	}

	return doc, err

}

func saveStringToFile(str string) (err error) {
	h := crc32.NewIEEE()
	h.Write([]byte(str))
	v := h.Sum32()
	bytes := strconv.Itoa(int(v))
	dirPath := DATA_DIR + bytes[0:2] + "/" + bytes[2:4]
	err = os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		return err
	}

	fileName := dirPath + "/" + bytes

	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(str)
	if err != nil {
		return err
	}

	return nil

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
	//	fmt.Println(links)
	//	getDoc("https://meduza.io/news/2018/03/11/v-kitae-otmenili-ogranichenie-sroka-prebyvaniya-u-vlasti-dlya-lidera-strany", m, &links)
	links.Add("https://meduza.io/news/2018/03/11/v-kitae-otmenili-ogranichenie-sroka-prebyvaniya-u-vlasti-dlya-lidera-strany")

	var i int = 0
	for next, exists := links.Remove(); exists; {
		fmt.Println("Links for downloading: ", len(links))
		fmt.Println("Total links: ", len(m))
		fmt.Println("Next: ", next)

		doc, err := goquery.NewDocument(next)
		if err != nil {
			fmt.Println("Error: ")
			fmt.Println(err)
		}

		stringForWrite, err := doc.Selection.Html()
		if err != nil {
			fmt.Println("Error: ")
			fmt.Println(err)
		}

		err = saveStringToFile(stringForWrite)
		if err != nil {
			fmt.Println("Error: ")
			fmt.Println(err)
		}

		//Example: get links from goquery doc
		getLinks(doc, m, &links)

		i++
		if i > MAX_LINKS {
			break
		}
		next, exists = links.Remove()

	}

	/*
		//Example: read goquery from file
		doc1, err := readDocumentFromFile(fileName)
		if err != nil {
			fmt.Println("Error: ")
			fmt.Println(err)
		}
	*/

	//fmt.Println(doc.Selection.Html())
}
