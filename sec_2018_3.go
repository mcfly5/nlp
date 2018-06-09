package main

import (
	"encoding/hex"
	"fmt"
)

func main() {
	src := []byte("19367831362e3d2b2c353d362c783136783336372f343d3c3f3d7839342f39212b782839212b782c303d783a3d2b2c7831362c3d2a3d2b2c")
	dst := make([]byte, hex.DecodedLen(len(src)))
	_, err := hex.Decode(dst, src)
	if err != nil {
		fmt.Println(err)
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
