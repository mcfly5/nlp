package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/agonopol/go-stem"
)

type Doc struct {
	length    int
	positions *[]int
}

type IndexNode struct {
	df   int
	docs map[string]*Doc
}

type kv struct {
	Key   string
	Value float64
}

var (
	DATA_DIR = "data/"
	//DATA_DIR = "old/"
	qntWords = 0
	qntDocs  = 0
	avgDocs  = 0
	//QUERY    = "QuickBooks Landlords"
	//QUERY = "Identity theft protection"
	index = make(map[string]*IndexNode)
)

func main() {

	//	buildIndex()
	buidIndexFromTsv("data/data.tsv")

	fmt.Println("Words :", qntWords)
	fmt.Println("Index:", len(index))
	loadQuery("data/topics.MB1-50.txt")

}

func topDocs(query string, N int) []string {

	candidates := make(map[string]float64)
	var result []string

	for _, word := range strings.Fields(query) {
		word = toStem(strings.ToLower(word))
		fmt.Println(word, ":", idf(word))
		if node, exist := index[word]; exist {
			for key, _ := range node.docs {
				candidates[key] = candidates[key] + tfidf(word, key)
				//fmt.Println("\t: key", key, " tfidf:", tf(word, key))
			}
			fmt.Println("\t:", len(node.docs))
		}
	}
	fmt.Println("Candidates:", len(candidates))

	var ss []kv
	for k, v := range candidates {
		ss = append(ss, kv{k, v})
	}

	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value > ss[j].Value
	})

	for i, kv := range ss {
		if i > N {
			return result
		}
		result = append(result, kv.Key)
		//fmt.Printf("i %d, %s, %d\n", i, kv.Key, kv.Value)
	}

	return result

	/*
		for key, val := range index {
			fmt.Println(key, ":")
			fmt.Println("\t", "Document frequency:", val.df)
			fmt.Println("\t", "Documents:")
			for path, doc := range val.docs {
				fmt.Println("\t\t", path, ":")
				fmt.Println("\t\t\t", "Length of document:", doc.length)
				fmt.Println("\t\t\t", "Term frequency:", len(*doc.positions), ",", *doc.positions)
			}

		}
	*/
	/*
		fmt.Println(stringToNGramm(QUERY))

		startTime := time.Now()

		if _, exist := dict[QUERY]; exist {
			fmt.Println("OK")
			return
		}

		minDistance := 100
		answer := ""
		maxFreq := 0
		i := 0
		var candidates []string
		grammes := stringToNGramm(QUERY)
		for _, gramm := range grammes {
			for word, _ := range index[gramm] {
				candidates = append(candidates, word)
			}
		}

		fmt.Println("Candidates:", len(candidates))

		for _, candidate := range candidates {
			distance := DamerauLevenshtein(candidate, QUERY)
			freq := dict[candidate]
			if distance < minDistance {
				answer = candidate
				minDistance = distance
				maxFreq = freq
			} else if distance == minDistance && maxFreq < freq {
				answer = candidate
				minDistance = distance
				maxFreq = freq
			}
			if i%1000 == 0 {
				fmt.Print(".")
			}
			i++
		}
		fmt.Println()
		fmt.Println("Candidate:", answer, "in distance:", minDistance)
		fmt.Println("Time:", time.Since(startTime))

	*/
}

func loadQuery(file string) {

	bs2, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println("Error while opening a file - ", file, "...")
		return
	}
	str := string(bs2)

	rows := strings.Split(str, "\n")
	for _, query := range rows {
		//if i > 9 { return		}
		val := strings.Split(query, ";")
		num, _ := strconv.Atoi(val[0])

		fmt.Print("i:", num)

		fmt.Println(" query:", val[1])

		fileout, err := os.Create("data/result/" + val[0] + ".txt")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer fileout.Close()

		for j, valOut := range topDocs(val[1], 99) {
			strOut := strconv.Itoa(j+1) + " " + valOut + "\n"
			_, err = fileout.WriteString(strOut)
			if err != nil {
				fmt.Println(err)
				return
			}

		}

	}
	fmt.Println(len(rows))

}

func idf(word string) float64 {

	var result float64 = 0

	if _, exist := index[word]; exist {
		result = math.Log10(float64(len(index)) / float64(index[word].df))
	}

	return result
}

func tf(word, doc string) float64 {

	var result float64 = 0

	if val, exist := index[word]; exist {
		if iDoc, exist := val.docs[doc]; exist {
			result = math.Log10(float64(1) + float64(len(*iDoc.positions))/float64(iDoc.length))
		}
	}

	return result
}

func bm25(word, doc string) float64 {
	//qntDocs - size of collection (N)
	//index[word].df - freq term(word) in collecion (DF)
	//len(index[word].docs[doc].positions) - freq term(word) in document(doc) (TF)
	//index[word].docs[doc].length length of document(doc)
	//avgDocs - average document lenght in collection
	var result float64 = 0

	//result = math.Log10

	return result

}

func tfidf(word, doc string) float64 {
	var result float64 = 0

	result = tf(word, doc) * idf(word)

	return result

}

func toStem(word string) string {

	return string(stemmer.Stem([]byte(word)))
}

func buildIndex() {

	filepath.Walk(DATA_DIR, buildIndexFromFiles)

	return
}

func buidIndexFromTsv(file string) {

	bs2, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println("Error while opening a file - ", file, "...")
		return
	}
	str := string(bs2)
	rows := strings.Split(str, "\n")
	//fmt.Println(len(rows))
	for _, row := range rows {
		//		if i > 9 {
		//			return
		//		}
		fields := strings.Split(row, "\t")
		//	fmt.Println("ID:", fields[0])
		//	fmt.Println("Text:", fields[1])
		addDocToIndex(fields[0], fields[1])
	}

}

func buildIndexFromFiles(path string, info os.FileInfo, err error) error {

	bs, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	str := string(bs)

	fmt.Println(path)

	addDocToIndex(path, str)

	return nil
}

func addDocToIndex(doc, str string) {

	words := strings.Fields(strings.ToLower(str))

	for i, word := range words {
		word = toStem(word)
		if _, exist := index[word]; !exist {
			newIndexNode := new(IndexNode)
			newIndexNode.df = 0
			newIndexNode.docs = make(map[string]*Doc)
			index[word] = newIndexNode
		}
		if _, exist := index[word].docs[doc]; !exist {
			index[word].df++
			index[word].docs[doc] = new(Doc)
			index[word].docs[doc].length = len(words)
			index[word].docs[doc].positions = new([]int)
		}
		*index[word].docs[doc].positions = append(*index[word].docs[doc].positions, i)
	}
	qntWords += len(words)
	qntDocs++

}

/*
for _, gramm := range nGramms {
	if _, exist := index[gramm]; !exist {
		index[gramm] = make(map[string]int)
	}
	if _, exist := index[gramm][word]; !exist {
		index[gramm][word] = 0
	}
	//index[gramm][word]++
*/
