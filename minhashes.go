package main

import (
	"fmt"
	"hash/crc32"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"strings"
)

type Doc struct {
	path      string
	shingles  map[string]int
	minhashes []uint32
}

var (
	DATA_DIR                     = "data/"
	shingles map[string][]uint32 = make(map[string][]uint32)
	docs     map[string]*Doc     = make(map[string]*Doc)
)

func main() {

	filepath.Walk(DATA_DIR, do1)

	fmt.Println(len(shingles))
	fmt.Println(len(docs))

	for key, val := range docs {
		fmt.Println(key) //, val, len(val))
		for key1, val1 := range docs {
			qnt := 0
			if key1 == key {
				continue
			}

			for i := 0; i < len(val.minhashes); i++ {
				if val.minhashes[i] == val1.minhashes[i] {
					qnt++
				}
			}

			//		fmt.Print("qnt: ", qnt," len_p",len(val), " len_c", len(val1), "; ")
			fmt.Println("\t", key1, " ", qnt, "/", 10, " JC:", float32(qnt)/float32(10))
		}
	}

}

func hash([]byte) uint {

	//fmt.Println(((i*2 + 1)*x+b) % p) % m))
	return 0
}

func do1(path string, info os.FileInfo, err error) error {
	fmt.Print(path)

	bs, err := ioutil.ReadFile(path)

	if err != nil {
		fmt.Println(err)
		return nil
	}

	str := strings.ToLower(string(bs))
	//fmt.Println(str)

	var qntOfHashes int = 10

	if _, exist := docs[path]; !exist {
		docs[path] = new(Doc)
		docs[path].shingles = make(map[string]int)
		docs[path].minhashes = make([]uint32, qntOfHashes)
		for i := 0; i < len(docs[path].minhashes); i++ {
			docs[path].minhashes[i] = uint32(math.MaxUint32)
		}

	}

	astr := strings.Fields(str)
	fmt.Print("; words:", len(astr))
	runes := []rune(str)

	used := make(map[string]bool)

	for i := 0; i < len(runes)-1; i++ {
		shingle := string(runes[i]) + string(runes[i+1])
		//fmt.Print(i, ":", runes[i], ":", shingle)
		if _, exist := used[shingle]; exist {
			//	fmt.Println(" - used")
			continue
		}

		if _, exist := shingles[shingle]; !exist {
			shingles[shingle] = hashes32([]byte(shingle), 422127559, qntOfHashes)
		}

		used[shingle] = true

		if _, exist := docs[path].shingles[shingle]; !exist {
			docs[path].shingles[shingle] = 0
			for i := 0; i < len(docs[path].minhashes); i++ {
				if docs[path].minhashes[i] > shingles[shingle][i] {
					docs[path].minhashes[i] = shingles[shingle][i]
				}
			}
		}
		//docs[path].shingles[shingle]++

	}

	fmt.Println("; shingles:", len(docs[path].shingles))

	fmt.Println(" ", docs[path])

	return nil
}

//fmt.Printf("Hashes: %v\n", hashes32([]byte("Hi"), 422127559, 10))
func hashes32(bytes []byte, seed uint32, qnt int) []uint32 {

	crc32q := crc32.MakeTable(0xD5828281)
	x := crc32.Checksum(bytes, crc32q)
	b := uint32(seed) //rand.Intn(math.MaxUint32)
	//w := uint(32)
	//M := uint(16)
	//hashCodeSizeDiff = w-M
	//c := (((i*2 + 1)*x+b)%p)%m
	//((hstart * (i*2 + 1)) + rand.Intn(bmax)) >>  hashCodeSizeDiff
	//	c := (a*x+b) >> (w-M)

	mins := make([]uint32, qnt)

	for i := 0; i < qnt; i++ {
		mins[i] = uint32(math.MaxUint32)
	}

	for i, min := range mins {
		t := (uint32(i*2)+1)*x + b
		//	fmt.Printf(" i: %#5d, CRC32: %#50b, %#20d\n", i, t, t)
		if t < min {
			mins[i] = t
		}
	}

	return mins
}

/*
func do(path string, info os.FileInfo, err error) error {
	fmt.Print(path)
	bs, err := ioutil.ReadFile(path)

	if err != nil {
		fmt.Println(err)
		return nil
	}

	str := strings.ToLower(string(bs))
	//fmt.Println(str)

	if _, exist := docs[path]; !exist {
		docs[path] = make(map[string]int)
	}

	astr := strings.Fields(str)
	fmt.Print("; words:", len(astr))

	runes := []rune(str)
	for i := 0; i < len(runes)-1; i++ {
		shingle := string(runes[i]) + string(runes[i+1])
		//fmt.Println(i, ":", runes[i], ":", shingle)
		if _, exist := shingles[shingle]; !exist {
			shingles[shingle] = 0
		}
		shingles[shingle]++

		if _, exist := docs[path][shingle]; !exist {
			docs[path][shingle] = 0
		}
		docs[path][shingle]++

	}
	fmt.Println("; shingles:", len(docs[path]))
	fmt.Println(docs[path])
	return nil
}
*/
