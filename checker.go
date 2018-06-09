package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	DATA_DIR = "data/imdb/neg/"
	//DATA_DIR = "old/"
	N        = 3
	qntWords = 0
	QUERY    = "mistuke"
	dict     = make(map[string]int)
	index    = make(map[string]map[string]int)
)

func main() {

	buildNGrammDict()

	fmt.Println("Words :", qntWords)
	fmt.Println("Uniqs :", len(dict))
	fmt.Println("Gramms:", len(index))

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
}

func stringToNGramm(str string) []string {

	bQuery := []byte(str)
	result := make([]string, len(bQuery)-N+1)
	for i := 0; i < len(bQuery)-N+1; i++ {
		result[i] = string(bQuery[i : i+N])
	}
	return result
}

func buildNGrammDict() {

	filepath.Walk(DATA_DIR, buildNGrammFromFile)

	return
}

func buildNGrammFromFile(path string, info os.FileInfo, err error) error {

	bs, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	str := string(bs)

	//fmt.Println(str)

	words := strings.Fields(str)

	for _, word := range words {
		if _, exist := dict[word]; !exist {
			if len(word) < N {
				word = word + strings.Repeat(" ", N-len(word))
			}
			nGramms := stringToNGramm(word)
			for _, gramm := range nGramms {
				if _, exist := index[gramm]; !exist {
					index[gramm] = make(map[string]int)
				}
				if _, exist := index[gramm][word]; !exist {
					index[gramm][word] = 0
				}
				//index[gramm][word]++
			}
			dict[word] = 0
		}
		dict[word]++

	}

	qntWords += len(words)

	return nil
}

func DamerauLevenshtein(s1 string, s2 string) (distance int) {
	// index by code point, not byte
	r1 := []rune(s1)
	r2 := []rune(s2)

	// the maximum possible distance
	inf := len(r1) + len(r2)

	// if one string is blank, we needs insertions
	// for all characters in the other one
	if len(r1) == 0 {
		return len(r2)
	}

	if len(r2) == 0 {
		return len(r1)
	}

	// construct the edit-tracking matrix
	matrix := make([][]int, len(r1))
	for i := range matrix {
		matrix[i] = make([]int, len(r2))
	}

	// seen characters
	seenRunes := make(map[rune]int)

	if r1[0] != r2[0] {
		matrix[0][0] = 1
	}

	seenRunes[r1[0]] = 0
	for i := 1; i < len(r1); i++ {
		deleteDist := matrix[i-1][0] + 1
		insertDist := (i+1)*1 + 1
		var matchDist int
		if r1[i] == r2[0] {
			matchDist = i
		} else {
			matchDist = i + 1
		}
		matrix[i][0] = min(min(deleteDist, insertDist), matchDist)
	}

	for j := 1; j < len(r2); j++ {
		deleteDist := (j + 1) * 2
		insertDist := matrix[0][j-1] + 1
		var matchDist int
		if r1[0] == r2[j] {
			matchDist = j
		} else {
			matchDist = j + 1
		}

		matrix[0][j] = min(min(deleteDist, insertDist), matchDist)
	}

	for i := 1; i < len(r1); i++ {
		var maxSrcMatchIndex int
		if r1[i] == r2[0] {
			maxSrcMatchIndex = 0
		} else {
			maxSrcMatchIndex = -1
		}

		for j := 1; j < len(r2); j++ {
			swapIndex, ok := seenRunes[r2[j]]
			jSwap := maxSrcMatchIndex
			deleteDist := matrix[i-1][j] + 1
			insertDist := matrix[i][j-1] + 1
			matchDist := matrix[i-1][j-1]
			if r1[i] != r2[j] {
				matchDist += 1
			} else {
				maxSrcMatchIndex = j
			}

			// for transpositions
			var swapDist int
			if ok && jSwap != -1 {
				iSwap := swapIndex
				var preSwapCost int
				if iSwap == 0 && jSwap == 0 {
					preSwapCost = 0
				} else {
					preSwapCost = matrix[maxI(0, iSwap-1)][maxI(0, jSwap-1)]
				}
				swapDist = i + j + preSwapCost - iSwap - jSwap - 1
			} else {
				swapDist = inf
			}
			matrix[i][j] = min(min(min(deleteDist, insertDist), matchDist), swapDist)
		}
		seenRunes[r1[i]] = i
	}

	return matrix[len(r1)-1][len(r2)-1]
}

// min of two integers
func min(a int, b int) (res int) {
	if a < b {
		res = a
	} else {
		res = b
	}

	return
}

// max of two integers
func maxI(a int, b int) (res int) {
	if a < b {
		res = b
	} else {
		res = a
	}

	return
}
