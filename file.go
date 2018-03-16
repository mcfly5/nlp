package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	fmt.Println("Hello")
	file, err := os.Open("test.txt")
	if err != nil {
		fmt.Println("Error while opening a file...")
		return
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		fmt.Println("Error while getting a filestat...")
		return
	}

	fmt.Printf("File size is %d \n", stat.Size())

	bs := make([]byte, stat.Size())
	_, err = file.Read(bs)
	if err != nil {
		fmt.Println("Error while reading a file...")
		return
	}

	str := string(bs)
	println("Content:", str)

	fmt.Println("Close a file...")

	file2, err := os.Create("test2.txt")
	if err != nil {
		println("Error while creating a file...")

	}
	defer file2.Close()

	file2.WriteString("Hello world!")

	bs2, err := ioutil.ReadFile("test2.txt")
	if err != nil {
		fmt.Println("Error while opening a file #2...")
		return
	}
	str2 := string(bs2)
	println("Content #2:", str2)

}
