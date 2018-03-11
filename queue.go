package main

import (
	"fmt"
)

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

func main() {
	fmt.Println("Hello")

	myq := make(Queue, 0)
	myq.Add("5")
	myq.Add("4")
	myq.Add("3")
	myq.Add("2")

	var i int = 0

	for next, ex := myq.Remove(); ex; {
		fmt.Println(next)
		fmt.Println(myq)
		i++
		if i > 7 {
			break
		}
		next, ex = myq.Remove()
	}

}
