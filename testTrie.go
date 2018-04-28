package main

import (
	"fmt"
	"os"
	"strings"
)

type Node struct {
	letters map[rune]*Node
	lemmes  []string
}

//func (Node) addLemme(str string) {}

var rootNode *Node

func main() {

	//Init tree
	rootNode = new(Node)
	rootNode.letters = make(map[rune]*Node)
	//fmt.Println(rootNode)

	str := "кот котенок котлета котик котофей комок комик ключ ключ ит"
	//str := "кот котенок комок"
	//str := "zzz"
	createTrie(str)

	//testTrie(str)
	//fmt.Println(rootNode)
	//	printTrie(0, rootNode, true)

	fmt.Println(rootNode)

	fmt.Println(findStringInTrie("кот"))
	fmt.Println(findStringInTrie("теленок"))
	fmt.Println(findStringInTrie("коти"))

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

	fmt.Println(rootNode)

	currentNode := rootNode
	//		fmt.Println("Root length: ", len(currentNode.letters))
	runes := []rune(word)

	//Reverse
	for i := len(runes) - 1; i >= 0; i-- {
		//Direct
		//for i := 0; i < len(runes); i++ {

		runeT := runes[i]
		fmt.Print("\t", "I:", i)
		fmt.Print("\t", "Rune: ", runeT)

		if _, ok := currentNode.letters[runeT]; !ok {
			//currentNode.lemmes = append(currentNode.lemmes, word)
			currentNode.letters[runeT] = new(Node)
			currentNode.letters[runeT].letters = make(map[rune]*Node)
			//fmt.Printf("\t#v+", currentNode.lemmes)
		}
		//fmt.Println()
		currentNode = currentNode.letters[runeT]
		//fmt.Println("\t", "RN:", rootNode)
		//fmt.Println("\t", "CN:", currentNode)
	}
	currentNode.lemmes = append(currentNode.lemmes, word)
	fmt.Print("\t", &currentNode.lemmes)
	fmt.Println("\t", "CN:", currentNode)

}

func findStringInTrie(word string) ([]string, string, bool) {
	println("Find word: ", word, ":")

	currentNode := rootNode
	//		fmt.Println("Root length: ", len(currentNode.letters))
	runes := []rune(word)

	//Reverse
	for i := len(runes) - 1; i >= 0; i-- {

		runeT := runes[i]
		//			fmt.Print("\t", "I:", i)
		//			fmt.Print("\t", "Rune: ", runeT)
		_, ok := currentNode.letters[runeT]
		if !ok {
			fmt.Print("\t", i)
			return currentNode.lemmes, string(runes[i+1 : len(runes)]), false
		}
		//fmt.Println()
		currentNode = currentNode.letters[runeT]
		//			fmt.Println("\t", "RN:", rootNode)
		//			fmt.Println("\t", "CN:", currentNode)
	}

	return currentNode.lemmes, "", true

}

/*
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
			fmt.Println(strings.Repeat("\t", level), string(key), " lemmes:", len(val.lemmes))
		}
		fmt.Println(strings.Repeat("\t", level), key)
		printTrie(level+1, val, mode)
	}
}

*/
func writeLog(str string) {
	file, err := os.OpenFile("log.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	file.WriteString(str)
	file.WriteString("\n")
}
