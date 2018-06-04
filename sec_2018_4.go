package main

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
)

func main() {

	src := []byte(*readStringFromFile("sec2018/task14.input"))
	//fmt.Println(src)
	//fmt.Println(string(src))
	i := 0
	prev := 0
	for i < len(src) {

		if src[i] == 10 {
			fmt.Println("New string:")
			//			fmt.Println(src[prev : i-1])
			decode(src[prev : i-1])

			i++
			prev = i + 1
			continue
		}
		/*
			fmt.Print(src[i], src[i+1])
			b1, ok := fromHexChar(src[i])
			if !ok {
				return
			}
			b2, ok := fromHexChar(src[i+1])
			if !ok {
				return
			}
			fmt.Printf("\t %d %d %#16b %d %v\n", b1, b2, uint(b1)<<4|uint(b2), uint(b1)<<4|uint(b2), string(uint(b1)<<4|uint(b2)))
		*/
		i += 2
	}

}

func decode(src []byte) {
	dst := make([]byte, hex.DecodedLen(len(src)))
	_, err := hex.Decode(dst, src)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(dst)

	freq := make(map[byte]int)

	for _, val := range dst {
		if _, exists := freq[val]; !exists {
			freq[val] = 0
		}
		freq[val] += 1
	}

	fmt.Println("Len:", len(dst))
	wordsMax := len(dst) / 5
	fmt.Println("Words max:", wordsMax)
	wordsMin := len(dst) / 8
	fmt.Println("Words min:", wordsMin)
	fmt.Println(freq)

	for key, val := range freq {
		if val > wordsMin && val < wordsMax {
			fmt.Println(key, val)
			candidate := key ^ (byte(32))
			fmt.Print("Candidate: ", candidate, " - ")
			for _, b := range dst {
				fmt.Print(string(b ^ candidate))
			}
			fmt.Println()
		}
	}
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

// fromHexChar converts a hex character into its value and a success flag.

func fromHexChar(c byte) (byte, bool) {

	switch {

	case '0' <= c && c <= '9':

		return c - '0', true

	case 'a' <= c && c <= 'f':

		return c - 'a' + 10, true

	case 'A' <= c && c <= 'F':

		return c - 'A' + 10, true

	}

	return 0, false

}
