package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type XMLDictionary struct {
	XMLName xml.Name `xml:"dictionary"`
	Lemmata Lemmata  `xml:"lemmata"`
}

type Lemmata struct {
	XMLName  xml.Name `xml:"lemmata"`
	LemmList []Lemma  `xml:"lemma"`
}

type Lemma struct {
	XMLName xml.Name `xml:"lemma"`
	Id      string   `xml:"id,attr"`
	Rev     string   `xml:"rev,attr"`
	L       L        `xml:"l"`
	F       []F      `xml:"f"`
}

type L struct {
	T string `xml:"t,attr"`
	G []G    `xml:"g"`
}

type G struct {
	V string `xml:"v,attr"`
}

type F struct {
	T string `xml:"t,attr"`
}

type Node struct {
	letters map[rune]*Node
	lemmes  []int
}

//func (Node) addLemme(str string) {}

var rootNode *Node

func main() {
	//Read dictionary
	rawXmlData := readStringFromFile("dict.opcorpora.xml")
	//fmt.Printf(rawXmlData)
	var data XMLDictionary
	xml.Unmarshal([]byte(*rawXmlData), &data)

	fmt.Println("Lemmes in XML  :", len(data.Lemmata.LemmList))
	fmt.Println("raw in XML  :", len(*rawXmlData))

	//Init tree
	rootNode = new(Node)
	rootNode.letters = make(map[rune]*Node)
	//fmt.Println(rootNode)

	//str := "кот котенок котлета котик котофей комок комик ключ ключ ти"
	//str := "кот котенок комок"
	//str := "zzz"
	//createTrie(&str)

	createTrieFromXML(&data)

	//testTrie(str)
	//fmt.Println(rootNode)
	//printTrie(0, rootNode, true)

	//fmt.Println(rootNode)

	fmt.Println(findStringInTrie("кот"))
	fmt.Println(findStringInTrie("подчеркивалось"))
	fmt.Println(findStringInTrie("отчужденно"))

}

func createTrieFromXML(data *XMLDictionary) {
	for _, lemma := range data.Lemmata.LemmList {
		id, err := strconv.Atoi(lemma.Id)
		if err != nil {
			fmt.Println(err)
		}
		for _, val := range lemma.F {
			addStringToTrie(id, val.T)
			//		writeLog(word)
		}

	}

}

func createTrie(str *string) {

	words := strings.Fields(*str)
	for _, word := range words {
		addStringToTrie(0, word)
		writeLog(word)
	}
}

func addStringToTrie(id int, word string) {
	//println("Add word: ", word, ":")

	//fmt.Println(rootNode)

	currentNode := rootNode
	//		fmt.Println("Root length: ", len(currentNode.letters))
	runes := []rune(word)

	//Reverse
	for i := len(runes) - 1; i >= 0; i-- {
		//Direct
		//for i := 0; i < len(runes); i++ {

		runeT := runes[i]
		//fmt.Print("\t", "I:", i)
		//fmt.Print("\t", "Rune: ", runeT)

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
	currentNode.lemmes = append(currentNode.lemmes, id)
	//fmt.Print("\t", &currentNode.lemmes)
	//fmt.Println("\t", "CN:", currentNode)

}

func findStringInTrie(word string) ([]int, string, bool) {

	var notFound bool = false
	var suffix string = ""

	println("Find word: ", word, ":")

	currentNode := rootNode
	//		fmt.Println("Root length: ", len(currentNode.letters))
	runes := []rune(word)

	//Reverse
	for i := len(runes) - 1; i >= 0; i-- {

		runeT := runes[i]
		//			fmt.Print("\t", "I:", i)
		//			fmt.Print("\t", "Rune: ", runeT)
		//fmt.Print("\t", "len(letters): ", len(currentNode.letters))
		_, ok := currentNode.letters[runeT]
		if !ok {
			fmt.Println("\t", i)
			notFound = true
			diff := ""
			for len(currentNode.letters) > 0 {
				for nextLetter, v := range currentNode.letters {
					currentNode = v
					diff = string(nextLetter) + diff
					break
				}
			}
			suffix = string(runes[i+1 : len(runes)])
			fmt.Println("\tOrig  : ", string(runes[0:i+1]))
			fmt.Println("\tDiff  : ", diff)
			fmt.Println("\tCommon: ", suffix)
			return currentNode.lemmes, suffix, notFound
		}
		//fmt.Println()
		currentNode = currentNode.letters[runeT]
		//			fmt.Println("\t", "RN:", rootNode)
		//			fmt.Println("\t", "CN:", currentNode)
	}

	return currentNode.lemmes, suffix, notFound

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
*/

func printTrie(level int, node *Node, mode bool) {
	for key, val := range node.letters {
		if mode {
			fmt.Println(strings.Repeat("\t", level), string(key), " lemmes:", len(val.lemmes))
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
	file.WriteString("\n")
}

func readStringFromFile(file string) *string {

	//	bs2, err := ioutil.ReadFile("dict_test.xml")

	str := new(string)

	bs2, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println("Error while opening a file - ", file, "...")
		return nil
	}
	*str = string(bs2)

	return str
}
