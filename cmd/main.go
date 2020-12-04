package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/kameshsampath/goignorescanner/pkg/scanner"
)

func main() {

	fmt.Println("Jai Guru")

	wd, err := os.Getwd()
	if err != nil {
		log.Fatal("os.Getwd() =", err)
	}

	//dir := filepath.Join(wd, "pkg", "scanner", "testdata", "dir2")
	dir := filepath.Join(wd, "pkg", "scanner", "testdata", "starignore")

	di := filepath.Join(dir, ".dockerignore")

	_, err = os.Stat(di)

	if err != nil {
		log.Fatal("os.Stat =", err)
	}

	ignoreScanner, err := scanner.NewOrDefault(dir)

	if err != nil {
		log.Fatal("scanner.NewOrDefault =", err)
	}

	includes, err := ignoreScanner.Scan()

	if err != nil {
		log.Fatal("ignoreScanner.Scan() =", err)
	}

	//ips := scanner.IgnorePatterns
	//
	//for _, e := range ips {
	//	fmt.Printf("\n Pattern file = %s \n", e.RegexPattern)
	//}

	var buf []byte

	buf = append(buf, "\033[22m"...)

	for _, incl := range includes {
		buf = append(buf, "\033[1m"...)
		buf = append(buf, "\n"...)
		buf = append(buf, incl...)
		buf = append(buf, "\033[1m"...)
	}

	_, err = os.Stdout.Write(append(buf, '\n'))

	if err != nil {
		log.Fatal(" os.Stdout.Write =", err)
	}
}
