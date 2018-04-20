package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

var (
	DATA_DIR    string = "data/"
	OUT_DIR     string
	RESCAN_TIME int = 10
)

func inquire(dirString string) error {
	dir, err := os.Open(dirString)
	if err != nil {
		return err
	}

	defer dir.Close()

	fileInfos, err := dir.Readdir(-1)
	if err != nil {
		return err
	}

	for i, fileInfo := range fileInfos {
		fmt.Println(i, ", Filename:", fileInfo.Name(), ", Size:", fileInfo.Size())
		fmt.Println("Content:")

		bs, err := ioutil.ReadFile(DATA_DIR + fileInfo.Name())
		if err != nil {
			return err
		}

		stringFromFile := string(bs)

		strSlice := strings.Split(stringFromFile, "\n")
		fmt.Println("Strings in file:", len(strSlice))
		for j, row := range strSlice {
			fmt.Println("\t Row#:", j)
			fmt.Println("\t", row)
			r := strings.NewReplacer(",", "", "-", "", ":", "", ";", "")
			newString := r.Replace(row)
			fmt.Println("\t New row:")
			fmt.Println("\t", newString)
		}

	}

	return nil

}

func payload() {
	for {
		fmt.Println("Do smth")
		time.Sleep(10 * time.Second)
	}
}

func f(n int) {
	for index := 0; index < 10; index++ {
		fmt.Println(n, ":", index)
		//amt := time.Duration(rand.Intn(250))
		time.Sleep(time.Millisecond * 1000)
	}
	//	time.Sleep(5 * time.Second)
}

func main2() {
	//	for i := 0; i < 10; i++ {
	//go f(i)
	go f(0)
	go payload()
	//	}
	var input string
	fmt.Scanln(&input)
	fmt.Println(input)

}

func main() {

	err := inquire(DATA_DIR)
	if err != nil {
		fmt.Println(err)
	}

	go payload()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, syscall.SIGTERM)
	fmt.Println("Press CTRL-C for exit")
	var i int = 0
	for {
		i++
		select {
		case <-sigChan:
			fmt.Println("Exiting")
			return
		}

		fmt.Println(i)
		time.Sleep(10 * time.Second)

	}

}
