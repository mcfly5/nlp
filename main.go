package main

import (
	"flag"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"strings"
	"time"
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

func getDoc(link string, m map[string]int, links *Queue) {
	doc, err := goquery.NewDocument(link)
	if err != nil {
		fmt.Println("Error: ")
		fmt.Println(err)
	}

	var max, maxValue int = 0, 0

	doc.Find("div").Each(func(i int, s *goquery.Selection) {

		length := len(s.Text())
		count := strings.Count(s.Text(), ".")
		if count > 2 && length/count < 250 {
			fmt.Println(s.Text())
		}
		if count > 2 && (length/count < maxValue) {
			max = i
			maxValue = length / count
		}

		//fmt.Printf("%d - %d - %s \n", length, count, s.Text())
	})
	fmt.Println(max, maxValue)

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
