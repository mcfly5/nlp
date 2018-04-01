package main

import (
	"fmt"
	"os"
	"strings"
)

type Node struct {
	letters map[rune]Node
	//lemmes  []int
}

var rootNode, currentNode Node

func main() {

	rootNode.letters = make(map[rune]Node)
	//fmt.Println(rootNode)

	str := "кот котенок котлета котик котофей комок комик"
	createTrie(str)

	//testTrie(str)
	//fmt.Println(rootNode)
	printTrie(0, rootNode, true)
}

func createTrie(str string) {

	words := strings.Fields(str)
	for _, word := range words {
		addStringToTrie(word)
		writeLog(word)
	}
}

func addStringToTrie(word string) {
	//println("Add word: ", word, ":")

	currentNode = rootNode
	//		fmt.Println("Root length: ", len(currentNode.letters))
	runes := []rune(word)

	//Reverse
	for i := len(runes) - 1; i >= 0; i-- {
		//Direct
		//for i := 0; i < len(runes); i++ {

		runeT := runes[i]
		//			fmt.Print("\t", "I:", i)
		//			fmt.Print("\t", "Rune: ", runeT)
		_, ok := currentNode.letters[runeT]
		if !ok {
			currentNode.letters[runeT] = Node{letters: make(map[rune]Node)}
			//				fmt.Print("\t", "New node!")
		}
		//fmt.Println()
		currentNode = currentNode.letters[runeT]
		//			fmt.Println("\t", "RN:", rootNode)
		//			fmt.Println("\t", "CN:", currentNode)
	}

}

func testTrie(str string) {

	var runeT rune

	words := strings.Fields(str)
	for _, word := range words {
		println(word, ":")
		currentNode = rootNode

		for _, runeT = range []rune(word) {
			_, ok := currentNode.letters[runeT]
			if !ok {
				currentNode.letters[runeT] = Node{letters: make(map[rune]Node)}
			}
			currentNode = currentNode.letters[runeT]
			//fmt.Println(rootNode)
			//fmt.Println(currentNode)
		}
	}

}

func printTrie(level int, node Node, mode bool) {
	for key, val := range node.letters {
		if mode {
			fmt.Println(strings.Repeat("\t", level), string(key))
		}
		fmt.Println(strings.Repeat("\t", level), key)
		printTrie(level+1, val, mode)
	}
}

func writeLog(str string) {
	file, err := os.OpenFile("log.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	file.WriteString(str)
}
